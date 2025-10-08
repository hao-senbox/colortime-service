package main

import (
	"colortime-service/config"
	"colortime-service/internal/colortime"
	"colortime-service/internal/language"
	"colortime-service/internal/product"
	"colortime-service/pkg/consul"
	"colortime-service/pkg/zap"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.LoadConfig()

	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	productService := product.NewUserService(consulClient)
	languageService := language.NewLanguageService(consulClient)
	colorTimeCollection := mongoClient.Database(cfg.MongoDB).Collection("colortime")

	colorTimeRepository := colortime.NewColorTimeRepository(colorTimeCollection)
	colorTimeService := colortime.NewColorTimeService(colorTimeRepository, productService, languageService)
	colorTimeHandler := colortime.NewColorTimeHandler(colorTimeService)
	router := gin.Default()

	colortime.RegisterRoutes(router, colorTimeHandler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Infof("Server running on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()

	// ✅ Graceful shutdown: chờ tín hiệu kill
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down server: %v", err)
	}
	logger.Info("Server stopped")
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}
