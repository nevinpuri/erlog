package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"erlog.net/netlogger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := zerolog.New(netlogger.NetLogger{}).With().Logger()
	// log := zerolog.New(netlogger.NetLogger{}).With().Logger()

	go func() {
		time.Sleep(time.Second * 2)
		logger.Print("log this")
		logger.Print("hello world")
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
		}

		fmt.Printf("%s\n", string(data))
	})

	// get all logs
	r.POST("/logs", func(c *gin.Context) {
	})

	log.Print(r.Run())
}