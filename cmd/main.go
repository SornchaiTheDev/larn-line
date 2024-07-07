package main

import (
	"larn-line/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := gin.Default()

	r.GET("/home", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	app, err := services.NewLineService(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	r.POST("/", app.Callback)

	r.Run(":3000")

}
