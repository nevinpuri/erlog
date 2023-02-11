package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/rs/zerolog/log"

	"erlog/models"
	"erlog/queue"

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
	// log.Logger = zerolog.New(netlogger.New()).With().Logger()

	go queue.Run()

	// todo: make queue take channel and return OK when it's started running
	// so we can block here and wait for the queue to start

	
	/*
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("Starting")
		for i := 0; i < 4; i++ {
			log.Print("good log")
		}
	}()
	*/
	

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)


		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// fmt.Printf("%#v", string(buffer.Bytes()))

		err = queue.Append(buffer.Bytes())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Status(200)
	})

	type LogGetRequestBody struct {
		Search string `json:"search"`
	}

	// get all logs
	r.POST("/search/logs", func(c *gin.Context) {
		var requestBody LogGetRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		fmt.Printf("%s\n", requestBody.Search)

		// var logs []models.ErLog
		logs := make([]models.ErLog, 0)
		models.DB.Where("data LIKE ?", "%" + requestBody.Search + "%").Find(&logs)

		// fmt.Printf("%s\n", logs[0].Data)

		// for i, log := range logs {
		// 	err := fastjson.ValidateBytes(log.O)

		// 	if err != nil {
		// 		fmt.Printf("%d:%s\n", i, err.Error())
		// 		c.JSON(http.StatusInternalServerError, err.Error())
		// 		return
		// 	}
		// }

		fmt.Printf("%d", len(logs))

		c.JSON(http.StatusOK, logs)
	})

	log.Print(r.Run())
}
