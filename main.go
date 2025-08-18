package main

import (
	"context"
	"demoService/src/application/contracts"
	"demoService/src/application/services"
	"demoService/src/infrastructure/cache"
	"demoService/src/infrastructure/database"
	"demoService/src/infrastructure/database/entities"
	"demoService/src/infrastructure/database/repositories"
	"demoService/src/infrastructure/message_querry"
	"demoService/src/web/controllers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "demoService/docs"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// @title Demo Service API
// @version 1.0
// @description Демонстрационный сервис с Kafka, PostgreSQL, кешем

func main() {
	log.Println("Hi, beauty 😘")
	loadEnv()
	db := initDatabase()
	lruCache := initCache()
	orderService := initServices(db, lruCache)
	startKafkaConsumer(orderService)
	startHTTPServer(orderService)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func initDatabase() *gorm.DB {
	db := database.DatabaseConnection()
	if err := db.AutoMigrate(&entities.Order{}, &entities.Delivery{}, &entities.Payment{}, &entities.Item{}); err != nil {
		log.Fatal(err)
	}
	return db
}

func initCache() contracts.Cache {
	cacheCapacity, err := strconv.Atoi(os.Getenv("CACHE_CAPACITY"))
	if err != nil {
		log.Fatal(err)
	}
	if cacheCapacity <= 0 {
		log.Fatal("CACHE_CAPACITY must be greater than zero")
	}

	lruCache, err := cache.NewLRUCache(cacheCapacity)
	if err != nil {
		log.Fatal(err)
	}

	cachePreload, err := strconv.Atoi(os.Getenv("CACHE_PRELOAD"))
	if err != nil {
		log.Fatal(err)
	}
	if cachePreload < 0 || cachePreload > cacheCapacity {
		log.Fatal("CACHE_PRELOAD out of bounds")
	}

	return lruCache
}

func initServices(db *gorm.DB, cache contracts.Cache) *services.OrderService {
	repo := repositories.NewOrderRepository(*db)
	orderService := services.NewOrderService(repo, cache)

	ctx := context.Background()

	cachePreload, _ := strconv.Atoi(os.Getenv("CACHE_PRELOAD"))
	if err := orderService.PreloadOrdersInCache(ctx, cachePreload); err != nil {
		log.Println("Error preloading cache:", err)
	}

	return orderService
}

func startKafkaConsumer(orderService *services.OrderService) {
	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	topic := os.Getenv("KAFKA_TOPIC")
	dlqTopic := os.Getenv("DLQ_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	dlqWriter := kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: dlqTopic,
	}

	workerCount, err := strconv.Atoi(os.Getenv("KAFKA_WORKERS"))
	if err != nil || workerCount <= 0 {
		workerCount = 12
	}

	consumer := message_querry.NewOrderConsumer(*orderService, reader, &dlqWriter)

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutdown signal received")
		cancel()
	}()

	go consumer.StartConsumer(ctx, workerCount)
}

func startHTTPServer(orderService *services.OrderService) {
	controller := controllers.NewOrderController(*orderService)
	controller.UseHandlers()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8081"
	}

	http.Handle("/swagger/", httpSwagger.Handler())

	fs := http.FileServer(http.Dir("src/web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server Started on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
