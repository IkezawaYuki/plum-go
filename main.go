package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"plum/infrastructure"
	"plum/presentation"
	"plum/usecase"
)

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://hp-standard.jp",
			"http://127.0.0.1",
		},
		AllowOriginWithContextFunc: nil,
		AllowMethods:               []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:               []string{"Authentication", "Origin"},
		AllowCredentials:           true,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping",
		})
	})

	hubspot := infrastructure.NewHubspot()
	slack := infrastructure.NewSlack()
	chatgpt := infrastructure.NewChatGPT(os.Getenv("AOAI_TOKEN"))

	contactService := usecase.NewContactService(hubspot, slack, chatgpt)
	handler := presentation.NewHandler(*contactService)

	r.POST("/support/contact", handler.SupportContact)

	if err := r.Run("localhost:8080"); err != nil {
		fmt.Println(err)
	}
}
