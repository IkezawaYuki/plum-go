package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"plum/infrastructure"
	"plum/presentation"
	"plum/usecase"
	"syscall"
)

func main() {

	r := gin.Default()
	presentation.Logger.Info("start plum!!!")

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

	hubspot := infrastructure.NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	slack := infrastructure.NewSlack()
	chatgpt := infrastructure.NewChatGPT(os.Getenv("AOAI_TOKEN"))
	gmailService := infrastructure.NewGmailService()

	contactService := usecase.NewContactService(hubspot, slack, chatgpt, gmailService)
	handler := presentation.NewHandler(*contactService)

	r.POST("/support/contact", handler.SupportContact)
	r.POST("/mail/hubspot", handler.GmailToHubspot)

	//if err := r.Run("localhost:8080"); err != nil {
	//	fmt.Println(err)
	//}
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			presentation.Logger.Error("listen: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	presentation.Logger.Info("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		presentation.Logger.Error("Server forced to shutdown: ", err)
	}
	presentation.Logger.Info("Server exiting")
}
