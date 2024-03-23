package routes

import (
	"context"
	"erlog/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

var perHour = "select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfMonth(timestamp) as date, toHour(timestamp) as hour, toMinute(timestamp) as minute, COUNT(*) as count from metrics GROUP BY minute, hour, date, month, year ORDER BY year, month, date, hour, minute;"
var perDay = "select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfMonth(timestamp) as date, COUNT(*) as count from metrics GROUP BY date, month, year ORDER BY year, month, date;"

func Search(c *gin.Context) {
	var searchQuery SearchRequestBody
	err := c.BindJSON(&searchQuery)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := ExecSearch(searchQuery.Per, "asdf")

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, data)
}

func Report(c *gin.Context) {
	var body ReportRequestBody
	err := c.BindJSON(&body)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	db.Conn.Exec(context.Background(), "INSERT INTO metrics VALUES (generateUUIDv4(), ?, now('Africa/Abidjan'))", body.Name)
	c.JSON(200, gin.H{"status": "ok"})
}