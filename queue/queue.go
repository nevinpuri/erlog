package queue

import (
	"errors"
	"fmt"
	"time"

	"erlog/asynclist"
	"erlog/models"
	"erlog/parser"

	"github.com/valyala/fastjson"
)

type Queue struct{
	BatchSize int
	timeout	int
	ch chan models.ErLog
	closeCh chan bool
	pool fastjson.ParserPool
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
		ch: make(chan models.ErLog),
		closeCh: make(chan bool),
		logs: asynclist.New(batchSize),
		pool: fastjson.ParserPool{},
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
			// append here

			q.logs.Append(log)

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

	// todo: instead of using fastjson.parsebytes, use parserpool
	p := q.pool.Get()
	val, err := p.ParseBytes(log)

	if err != nil {
		return err
	}

	q.pool.Put(p)

	// ideally this should be an object so just get the object here but it really doesn't matter for now
	erlog, err := parser.ParseJson(val)

	if err != nil {
		return err
	}

	erlog.Raw = string(log)

	q.ch <- erlog

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

	batch, err := models.Conn.PrepareBatch(models.CTX, "INSERT INTO er_logs")

	if err != nil {
		return err
	}

	for _, v := range q.logs.All() {
		// fmt.Printf("%d", v.Id)

		fmt.Printf("String keys: %v, string values: %v, number keys: %v, number values: %v, bool keys: %v, bool_values: %v\n", v.StringKeys, v.StringValues, v.NumberKeys, v.NumberValues, v.BoolKeys, v.BoolValues)
		err = batch.AppendStruct(&v)

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	err = batch.Send()

	fmt.Printf("Send logs")

	if err != nil {
		fmt.Printf("%v", err)
		return err
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