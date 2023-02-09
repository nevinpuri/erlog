package netlogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NetLogger struct {
	queue []string
}

func New() NetLogger {
	return NetLogger{}
}

func (logger NetLogger) Write(p []byte) (n int, err error) {
	if !json.Valid(p) {
		fmt.Printf("%s is not valid, skipping\n", p)
		return
	}

	_, err = http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(p))

	if err != nil {
		return len(p), err
	}

	return len(p), nil
}