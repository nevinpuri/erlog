package queue

import (
	"errors"
	"fmt"
	"time"

	"erlog.net/asynclist"
	"erlog.net/models"
)

type Queue struct{
	BatchSize int
	timeout	int
	ch chan []byte
	closeCh chan bool
	timer *time.Timer
	logs asynclist.AsyncList
}

// creates a new queue
// batchSize: the size to batch items before flushing
// timeout: the period of time (in ms) to wait before flushing
func New(batchSize int, timeout int) Queue {
	return Queue {
		BatchSize: batchSize,
		// maybe make the size runtime.NumCPU() but idk
		ch: make(chan []byte),
		closeCh: make(chan bool),
		logs: asynclist.New(batchSize),
		timer: time.NewTimer(time.Millisecond * time.Duration(timeout)),
		timeout:timeout,
	}
}

// TODO: make flush after x seconds code
func (q *Queue) Run() {
	for {
		select {
		case <- q.closeCh:
			// according to go:
			// only the sender should close the channel,
			// never the reciever
			// close(q.ch)
			// close(q.closeCh)

			err := q.Flush()

			// don't know of much else we can do here, good luck
			if err != nil {
				// log.Fatalf("%s", err.Error())
				fmt.Printf("%s", err.Error())
			}

			return
		case log := <- q.ch:
			if len(log) == 0 {
				continue
			}

			q.logs.Append(log)

			// fmt.Printf("Got new log %s", log)

			if q.logs.Len() < q.BatchSize {
				continue
			}

			err := q.Flush()

			if err != nil {
				fmt.Printf("%s", err.Error())
			}

		case <- q.timer.C:
			// fmt.Printf("Flushing because of timer %d", q.timeout)
			err := q.Flush()

			if err != nil {
				fmt.Printf("%s", err.Error())
			}

			q.timer.Reset(time.Millisecond * time.Duration(q.timeout))
		}
	}
}

func (q *Queue) Append(log []byte) error {
	// validate the json here
	q.ch <- log

	return nil
}

func (q *Queue) Flush() error {

	logsLen := q.logs.Len()

	if logsLen == 0 {
		return nil
	}

	if !models.IsConnected() {
		return errors.New("Db not connected")
	}

	fmt.Println("Flushing logs")

	erlogs := make([]models.ErLog, logsLen + 1)

	for i, v := range q.logs.All() {
		erlogs[i] = models.ErLog{O: v}
	}

	// for _, v := range erlogs {
	// 	fmt.Printf("Erlog: %s\n", v.O)
	// }

	// for _, v := range q.logs.All() {
	// 	fmt.Printf("%s\n", v)
	// }

	models.DB.CreateInBatches(erlogs, q.BatchSize)

	q.logs.Clear()

	return nil
}

func (q *Queue) Close() error {
	err := q.Flush()

	if err != nil {
		return err
	}

	q.closeCh <- true
	return nil
}