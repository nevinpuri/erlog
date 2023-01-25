package queue

import (
	"fmt"
	"time"

	"erlog.net/asynclist"
)

type Queue struct{
	BatchSize int
	ch chan []byte
	closeCh chan bool
	timer time.Timer
	logs asynclist.AsyncList
}

func (q *Queue) New(batchSize int) Queue {
	return Queue {
		BatchSize: batchSize,
		// maybe make the size runtime.NumCPU() but idk
		ch: make(chan []byte),
		closeCh: make(chan bool),
		logs: asynclist.New(batchSize),
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

			if q.logs.Len() < q.BatchSize {
				continue
			}

			err := q.Flush()

			if err != nil {
				fmt.Printf("%s", err.Error())
			}

		// TODO
		// case <- q.timer.C:
		}
	}
}

func (q *Queue) Append(log []byte) error {
	// validate the json here
	q.ch <- log

	return nil
}

func (q *Queue) Flush() error {
	if q.logs.Len() == 0 {
		return nil
	}

	// flush the logs to sqlite

	fmt.Println("Flushing logs")
	for _, v := range q.logs.All() {
		fmt.Printf("%s\n", v)
	}

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