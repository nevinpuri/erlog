package netlogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var Wg sync.WaitGroup

type NetLogger struct {
	// queue	queue.Queue
}

func New() NetLogger {
	return NetLogger{}
	// return NetLogger{queue: queue.New(batchSize)}
}

func (logger NetLogger) Write(p []byte) (n int, err error) {

	if !json.Valid(p) {
		fmt.Printf("%s is not valid, skipping\n", p)
		// Wg.Done()
		return
	}

	http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(p))

	// fmt.Print(string(p))

	return len(p), nil
}