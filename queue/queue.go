package queue

import "log"

type Queue struct{
	BatchSize int
	ch chan []byte
	closeCh chan bool
	len int
}

func (q Queue) New(batchSize int) Queue {
	return Queue {
		BatchSize: batchSize,
		// maybe make the size runtime.NumCPU() but idk
		ch: make(chan []byte),
	}
}

// TODO: make flush after x seconds code
func (q Queue) Run() {
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
				log.Fatalf("%s", err.Error())
			}

			return
		}
	}
}

func (q Queue) Append(log []byte) error {
	// validate the json here
	q.len += 1
	q.ch <- log

	return nil
}

func (q Queue) Flush() error {
	if q.len == 0 {
		return nil
	}
	
	// batch the items and append them to sqlite

	q.len = 0
	// todo: append to database
	return nil
}

func (q Queue) Close() error {
	err := q.Flush()

	if err != nil {
		return err
	}

	q.closeCh <- true
	return nil
}