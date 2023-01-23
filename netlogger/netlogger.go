package netlogger

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

var Wg sync.WaitGroup

type NetLogger struct {
}

func (logger NetLogger) Write(p []byte) (n int, err error) {
	Wg.Add(1)
	go func() {
		http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(p))

		fmt.Print(string(p))
		Wg.Done()
	}()

	return 0, nil
}