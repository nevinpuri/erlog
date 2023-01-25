package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"erlog.net/netlogger"
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

	log.Logger = zerolog.New(netlogger.NetLogger{}).With().Logger()

	go func() {
		time.Sleep(time.Second * 2)
		log.Print("log this")
		log.Print("hello world")
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

		// result := json.RawMessage{}
		// err = json.Unmarshal(data, &result)

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error(),
		// 	})
		// }

		db.Create(&ErLog{O: data})

		fmt.Printf("%s\n", string(data))

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