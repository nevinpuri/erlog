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
	"github.com/valyala/fastjson"

	"erlog/converter"
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

	r.POST("/logs", func (c *gin.Context) {
		var logs []models.ErLog

		// todo: make this application wide (maybe)
		converter := converter.New()

		if err := models.Conn.Select(models.CTX, &logs, "SELECT * FROM er_logs"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var objs []fastjson.Object

		for _, log := range logs {
			obj, err := converter.Convert(log)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			objs = append(objs, obj)
		}

		fmt.Printf("%v\n", len(objs))

		c.JSON(http.StatusOK, objs)
	})

	log.Print(r.Run())
}
