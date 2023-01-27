package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

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

	// todo: make queue take channel and return OK when it's started running
	// so we can block here and wait for the queue to start

	/*
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
	*/

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

		err = queue.Append(data)

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
		models.DB.Find(&logs).Limit(2)

		fmt.Printf("%s\n", logs[0].O)

		// for i, log := range logs {
		// 	err := fastjson.ValidateBytes(log.O)

		// 	if err != nil {
		// 		fmt.Printf("%d:%s\n", i, err.Error())
		// 		c.JSON(http.StatusInternalServerError, err.Error())
		// 		return
		// 	}
		// }

		fmt.Printf("%d", len(logs))

		c.JSON(http.StatusOK, gin.H{"data": logs})
	})

	log.Print(r.Run())
}
