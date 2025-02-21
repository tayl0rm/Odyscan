package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// Load HTML templates from 'templates' folder
	r.LoadHTMLGlob("templates/*")

	// Serve static files from 'static' folder
	r.Static("/static", "./static")

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		})
	})

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start server
	r.Run(":8080")
}
