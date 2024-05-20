package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"plum/infrastructure"
	"plum/logger"
	"plum/presentation"
	"plum/usecase"
	"syscall"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	logger.Logger.Info("start plum!!!")

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

	r.Static("/plum/static", "./static")

	loadTemplates := func(router *gin.Engine) {
		templ := template.Must(template.New("").Funcs(template.FuncMap{}).ParseGlob("templates/*.tmpl"))
		r.SetHTMLTemplate(templ)
	}

	loadTemplates(r)

	r.GET("/plum/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping",
		})
	})

	hubspot := infrastructure.NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	slack := infrastructure.NewSlack(os.Getenv("SLACK_WEBHOOK_URL"))
	chatgpt := infrastructure.NewChatGPT(
		os.Getenv("AZURE_OPENAI_KEY"),
		os.Getenv("AZURE_OPENAI_ENDPOINT"))
	fmt.Println(os.Getenv("AZURE_OPENAI_KEY"))
	fmt.Println(os.Getenv("AZURE_OPENAI_ENDPOINT"))
	b, err := os.ReadFile("./credentials.json")
	if err != nil {
		log.Fatalf("%v", err)
	}
	gmailService := infrastructure.NewGmailService(b, "./token.json")
	aiSearchSearch := infrastructure.NewAISearch(
		os.Getenv("AI_SEARCH_URL"),
		os.Getenv("AI_SEARCH_API_KEY"))

	contactService := usecase.NewContactService(
		hubspot,
		slack,
		chatgpt,
		gmailService,
		aiSearchSearch,
	)
	handler := presentation.NewHandler(*contactService)

	r.POST("/plum/support/form", handler.SupportForm)
	r.POST("/plum/support/mail", handler.SupportMail)
	r.GET("/plum/support/form", handler.SupportFormPage)
	r.GET("/plum/support/thank_you", handler.ThankYouPage)

	server := &http.Server{
		Addr:    "0.0.0.0:8001",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Logger.Error("listen: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Logger.Error("Server forced to shutdown: ", err)
	}
	logger.Logger.Info("Server exiting")
}
