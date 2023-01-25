package main

import (
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

func main() {
	models.Connect()

	queue := queue.New(512, 3000)
	log.Logger = zerolog.New(netlogger.New()).With().Logger()

	go queue.Run()

	go func() {
		time.Sleep(time.Second * 2)

		fmt.Println("Starting")
		start := time.Now()
		for i := 0; i < 2400; i++ {
			log.Print("log this")
			log.Print("hello world")
			log.Print("new logs")
			log.Print("final")
		}

		elapsed := time.Since(start)

		fmt.Printf("%s", elapsed)

		fmt.Println("Done")

		time.Sleep(time.Second * 2)
		log.Print("HI HI HI")
		log.Print("Miami")
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
		var logs []models.ErLog
		models.DB.Find(&logs)

		c.JSON(http.StatusOK, logs)
	})

	log.Print(r.Run())
}