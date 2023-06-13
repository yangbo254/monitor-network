package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartHttpEngine() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/timeData", func(c *gin.Context) {
		now := time.Now().Unix()
		begin := now - 60*60*6
		data := DbGetVistorRecords(begin, now)
		c.JSON(http.StatusOK, data)
	})
	r.Static("/html", "./html")
	go r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
