package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"erlog.net/netlogger"
	"erlog.net/queue"
	"github.com/gin-gonic/gin"
)

type ErLog struct {
	gorm.Model
	O	datatypes.JSON
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		fmt.Printf("Failed opening db\nError:%s", err.Error());
		return
	}

	db.AutoMigrate(&ErLog{})

	queue := queue.New(2, 5000)
	log.Logger = zerolog.New(netlogger.New()).With().Logger()

	go queue.Run()

	go func() {
		log.Print("log this")
		log.Print("hello world")
		log.Print("new logs")
		log.Print("final")
	}()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		queue.Append(data)

		// result := json.RawMessage{}
		// err = json.Unmarshal(data, &result)

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error(),
		// 	})
		// }

		// db.Create(&ErLog{O: data})

		// fmt.Printf("%s\n", string(data))

		c.Status(200)
	})

	// get all logs
	r.GET("/logs", func(c *gin.Context) {
		var logs []ErLog
		db.Find(&logs)

		c.JSON(http.StatusOK, logs)
	})

	log.Print(r.Run())
}