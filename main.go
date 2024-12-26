package main

import (
	"context"
	"fmt"
	"log"
	"myproject/middlewares"
	"myproject/routers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var updatedVersion = "1.0.32"

func init() {
	fmt.Println("::Environment mode : " + gin.Mode())
	if gin.Mode() == "debug" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		fmt.Println("::Environment Variables : loaded from .env")
	}

	err := mgm.SetDefaultConfig(nil, "cognify", options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		log.Fatalf("::DB Connection Error : Failed to connect to MongoDB: %v", err)
	}

	fmt.Println("::DB Connection Status : Successfully connected to MongoDB!")
}

func main() {
	router := gin.Default()

	// Apply Middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.RateLimitingMiddleware())
	router.Use(middlewares.CORSMiddleware())
	//router.Use(middlewares.JWTAuthMiddleware())
	router.Use(middlewares.ErrorHandlingMiddleware())
	router.Use(middlewares.InputValidationMiddleware())

	// Routes
	router.POST("/auth", routers.Authenticate)
	router.POST("/questions", routers.SubmitQuestions)
	router.GET("/questions", routers.FetchAllQuestions)
	router.POST("/submit", routers.HandleSubmission)
	router.GET("/report", routers.HandleReportGeneration)
	router.GET("/paymentCallback", routers.HandlePaymentCallback)
	router.GET("/report/:testId", routers.HandleBig5Report)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy v:" + updatedVersion})
	})

	playgroundRouter := os.Getenv("PLAYGROUND_ROUTER")
	if playgroundRouter == "allowed" {
		//Pdf test route
		router.POST("/pdf", routers.CreatingPdf)
		router.POST("/mail", routers.TestMail)
		router.POST("/generatepdf", routers.Generatepdf)
		router.POST("/paymentLinkCreate", routers.PaymentTest)
		router.POST("/paymentLinkFetch", routers.PaymentLinkFetch)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback to port 8080 if not set
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
