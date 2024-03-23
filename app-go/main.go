package main

import (
	"erlog/db"
	"erlog/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func (c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/report", routes.Report)
	r.POST("/search", routes.Search)

	http.ListenAndServe(":8000", r)

}
