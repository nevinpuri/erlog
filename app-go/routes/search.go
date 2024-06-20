package routes

import (
	"context"
	"erlog/db"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)


var perHour = "select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfMonth(timestamp) as date, toHour(timestamp) as hour, toMinute(timestamp) as minute, COUNT(*) as count from metrics GROUP BY minute, hour, date, month, year ORDER BY year, month, date, hour, minute;"
var perDay = "select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfMonth(timestamp) as date, toHour(timestamp) as hour, COUNT(*) as count from metrics GROUP BY date, month, year, hour ORDER BY year, month, date, hour;"

func ExecSearch(per string, name string) ([]SearchResponse, error) {
	var query string
	if per == "hour" {
		query = perHour
	} else if per == "day" {
		query = perDay
	} else {
		fmt.Printf("Not implemented\n")
		return nil, fmt.Errorf("not implemented query")
	}

	result, err := db.Conn.Query(context.Background(), query)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return nil, err
	}

	data, err := GetDataFor(per, result)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return nil, err
	}

	return data, nil
}

// gets the data
func GetDataFor(per string, rows driver.Rows) ([]SearchResponse, error) {
	var data []SearchResponse
	var err error
	if per == "hour" {
		data, err = GetDataPerHour(rows, data)
	} else if per == "day" {
		data, err = GetDataPerDay(rows, data)
	} else {
		return nil, fmt.Errorf("not implemented")
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	return data, nil
}

// returns data in per minute intervals
func GetDataPerHour(rows driver.Rows, data []SearchResponse) ([]SearchResponse, error) {
	for rows.Next() {
		var year uint16
		var month uint8
		var date uint8
		var hour uint8
		var minute uint8
		var count uint64

		if err := rows.Scan(&year, &month, &date, &hour, &minute, &count); err != nil {
			fmt.Printf("%v\n", err.Error())
			return nil, err
		}

		concatted := fmt.Sprint(date) + " " + fmt.Sprint(hour) + ":" + fmt.Sprint(minute)
		fmt.Printf("%v:%v\n", concatted, count)
		// data[concatted] = count
		data = append(data, SearchResponse{
			DateTime: concatted,
			Count:    count,
		})
	}

	return data, nil
}

// returns data in per hour intervals
func GetDataPerDay(rows driver.Rows, data []SearchResponse) ([]SearchResponse, error) {
	for rows.Next() {
		var year uint16
		var month uint8
		var date uint8
		var hour uint8
		var count uint64

		if err := rows.Scan(&year, &month, &date, &hour, &count); err != nil {
			fmt.Printf("%v\n", err.Error())
			return nil, err
		}

		concatted := fmt.Sprint(year) + "-" + fmt.Sprint(month) + "-" + fmt.Sprint(date) + " " + fmt.Sprint(hour)
		fmt.Printf("%v:%v\n", concatted, count)
		data = append(data, SearchResponse{
			DateTime: concatted,
			Count:    count,
		})
	}

	return data, nil
}