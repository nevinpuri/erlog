package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	err := models.Connect()

	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	models.Conn.Exec(models.CTX, models.SetupTable)

	queue := queue.New(512, 3000)

	go queue.Run()

	// todo: make queue take channel and return OK when it's started running
	// so we can block here and wait for the queue to start

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

		buffer := new(bytes.Buffer)
		err = json.Compact(buffer, data)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = queue.Append(buffer.Bytes())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Status(200)
	})

	type SearchRequestBody struct {
		Query string `json:"query" binding:"required"`
	}

	r.POST("/search", func (c *gin.Context) {
		var requestBody SearchRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// rows, err := models.Conn.Query(models.CTX, "SELECT * FROM er_logs WHERE indexOf(StringNames, ?) > 0", requestBody.Query)

		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
	})

	log.Print(r.Run())
}
