package routes

import (
	"context"
	"erlog/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

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