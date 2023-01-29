package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"erlog.net/models"
	"erlog.net/netlogger"
	"erlog.net/queue"
	"github.com/gin-gonic/gin"
)

type Event struct {
	Message chan string
	NewClients chan chan string
	ClosedClients chan chan string
	TotalClients map[chan string]bool
}

type ClientChan chan string

func main() {
	models.Connect()

	queue := queue.New(512, 3000)
	log.Logger = zerolog.New(netlogger.New()).With().Logger()

	go queue.Run()

	// todo: make queue take channel and return OK when it's started running
	// so we can block here and wait for the queue to start

	
	// /*
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("Starting")
		log.Print("log this")
	}()
	// */
	

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

		// trimmed := strings.TrimSpace(string(data))

		// c.JSON(http.StatusOK, trimmed)
		// return

		buffer := new(bytes.Buffer)
		err = json.Compact(buffer, data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// fmt.Printf("%#v", string(buffer.Bytes()))

		err = queue.Append(buffer.Bytes())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Status(200)
	})

	// get all logs
	r.GET("/logs", func(c *gin.Context) {
		var logs []models.ErLog
		models.DB.Find(&logs)

		fmt.Printf("%s\n", logs[0].Data)

		// for i, log := range logs {
		// 	err := fastjson.ValidateBytes(log.O)

		// 	if err != nil {
		// 		fmt.Printf("%d:%s\n", i, err.Error())
		// 		c.JSON(http.StatusInternalServerError, err.Error())
		// 		return
		// 	}
		// }

		fmt.Printf("%d", len(logs))

		c.JSON(http.StatusOK, logs[0])
	})

	log.Print(r.Run())
}
