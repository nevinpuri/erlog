package db

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var Conn driver.Conn

func ConnectDB() {
	dialCount := 0
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
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(0)
	}

	conn.Ping(context.Background())

	// so will basically be
	// metrics:
	// id, name, timestamp
	// report("user signed up")

	fmt.Println("Creating Metrics table")
	err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS metrics (id UUID, name String, timestamp DateTime) Engine = MergeTree PRIMARY KEY (id, timestamp, name) ORDER BY (id, timestamp, name);")
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(0)
	}
	fmt.Println("Done creating Metrics Table")

	Conn = conn
}