package main

import (
	"log"
	"os"
	"strings"

	"war-ticket-engine/handlers"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize services
	redisService := services.NewRedisService()
	captchaService := services.NewCaptchaService()

	// Initialize JSON Database
	services.InitDatabase("database.json")

	// Initialize RAG service
	ragPath := os.Getenv("RAG_DOC_PATH")
	if ragPath == "" {
		ragPath = "../README.md"
	}
	ragService, err := services.NewRAGService(ragPath)
	if err != nil {
		log.Printf("RAG disabled: %v", err)
	}

	// Initialize router
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	api := r.Group("/api")
	{
		// War tiket (original)
		api.POST("/war", handlers.WarHandler(redisService))
		api.GET("/status", handlers.StatusHandler(redisService))

		// Authentication
		api.POST("/register", handlers.RegisterHandler(captchaService))
		api.POST("/login", handlers.LoginHandler(captchaService))

		// Captcha
		api.GET("/captcha/math", handlers.MathCaptchaHandler(captchaService))
		api.GET("/captcha/image", handlers.ImageCaptchaHandler(captchaService))

		// Tickets & Locations
		api.POST("/ticket", handlers.CreateTicketHandler())
		api.GET("/ticket/:id", handlers.GetTicketHandler())
		api.GET("/locations", handlers.GetLocationsHandler())

		if ragService != nil {
			api.POST("/chat", handlers.ChatHandler(ragService))
		}
	}

	log.Println("Server starting")
	log.Println("Database: database.json")
	log.Println("Endpoints available")
	addr := strings.TrimSpace(os.Getenv("SERVER_ADDR"))
	if addr == "" {
		port := strings.TrimSpace(os.Getenv("PORT"))
		if port == "" {
			port = "8080"
		}
		if strings.HasPrefix(port, ":") {
			addr = port
		} else {
			addr = ":" + port
		}
	}
	log.Printf("Listening on %s", addr)
	r.Run(addr)
}
