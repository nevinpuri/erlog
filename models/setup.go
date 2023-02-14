package models

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var CTX context.Context
var Conn clickhouse.Conn
var dialCount int32
// var DB *gorm.DB

func Connect() {
	CTX = context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:19000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "test123",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...interface{}) {
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout: time.Duration(10) * time.Second,
		MaxOpenConns: 5,
		MaxIdleConns: 5,
		ConnMaxLifetime: time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		BlockBufferSize: 10,
	})

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	Conn = conn
}

// todo: make this actually check if the db is connected
func IsConnected() bool {
	if Conn == nil {
		return false
	}

	return true
}