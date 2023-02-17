package queue

import (
	"errors"
	"fmt"
	"time"

	"erlog/asynclist"
	"erlog/models"
	"erlog/parser"

	"github.com/google/uuid"
	"github.com/valyala/fastjson"
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
				fmt.Printf("%s", err.Error())
			}

			return
		case log := <- q.ch:
			if len(log) == 0 {
				continue
			}

			// append here

			// q.logs.Append(log)

			if q.logs.Len() < q.BatchSize {
				continue
			}

			err := q.Flush()

			if err != nil {
				fmt.Printf("%s", err.Error())
			}

		case <- q.timer.C:
			err := q.Flush()

			if err != nil {
				fmt.Printf("%s", err.Error())
			}

			q.timer.Reset(time.Millisecond * time.Duration(q.timeout))
		}
	}
}

func (q *Queue) Append(log []byte) error {
	// var js interface{}

	// err := json.Unmarshal(log, &js)

	// if err != nil {
	// 	return err
	// }

	val, err := fastjson.ParseBytes(log)

	if err != nil {
		return err
	}

	keys := ""
	erlog := models.ErLog{}

	// ideally this should be an object so just get the object here but it really doesn't matter for now
	parser.ParseValue(val, &keys, &erlog)

	// we don't actually care about the value of js
	// we just care that it's valid
	
	q.ch <- log

	return nil
}

func (q *Queue) Flush() error {

	logsLen := q.logs.Len()

	if logsLen == -1 {
		return nil
	}

	if !models.IsConnected() {
		return errors.New("Db not connected")
	}

	erlogs := make([]models.ErLog, logsLen + 1)

	for i, v := range q.logs.All() {
		fmt.Printf("%d", v.Id)
		// parse the fucking log

		erlogs[i] = models.ErLog{
			Id: uuid.New(),
		}
	}

	// batch, err := models.Conn.PrepareBatch(models.CTX, "INSERT INTO er_logs")

	// if err != nil {
	// 	return err
	// }

	// models.DB.PrepareBatch()
	// models.DB.CreateInBatches(erlogs, q.BatchSize)

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