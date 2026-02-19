package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/gino/cars-crud/docs"

	"github.com/gino/cars-crud/internal/cache"
	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/internal/handler"
	"github.com/gino/cars-crud/internal/middleware"
	"github.com/gino/cars-crud/internal/queue"
	mongoRepo "github.com/gino/cars-crud/internal/repository/mongo"
	pgRepo "github.com/gino/cars-crud/internal/repository/postgres"
	"github.com/gino/cars-crud/internal/usecase"
	"github.com/gino/cars-crud/pkg/config"
)

// @title        Cars CRUD API
// @version      1.0
// @description  A RESTful API for managing cars with Go, Chi, GORM, Redis, Kafka, and MongoDB.

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token. Example: "Bearer eyJhbG..."

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	if err := db.AutoMigrate(&domain.Car{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println("postgres connected and migrated")

	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("redis connected")

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	log.Println("mongo connected")

	producer := queue.NewLogProducer(cfg)
	defer producer.Close()

	consumer := queue.NewLogConsumer(cfg, mongoClient)
	consumer.Start(ctx)
	defer consumer.Close()
	log.Println("kafka producer and consumer started")

	logCollection := mongoClient.Database(cfg.MongoDB).Collection(cfg.MongoCollection)

	carRepo := pgRepo.NewCarRepository(db)
	logRepo := mongoRepo.NewLogRepository(logCollection)
	carUsecase := usecase.NewCarUsecase(carRepo, redisCache)

	authHandler := handler.NewAuthHandler(cfg.APIKey, cfg.JWTSecret)
	carHandler := handler.NewCarHandler(carUsecase)
	logHandler := handler.NewLogHandler(logRepo)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestLogger(producer))

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	authHandler.RegisterRoutes(r)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		carHandler.RegisterRoutes(r)
		logHandler.RegisterRoutes(r)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
