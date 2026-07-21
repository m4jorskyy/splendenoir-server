package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"splendenoir-server/internal/handlers"
	"splendenoir-server/internal/middleware"
	"splendenoir-server/internal/repositories"
	"splendenoir-server/internal/services"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var waitGroup sync.WaitGroup

	errEnv := godotenv.Load()
	if errEnv != nil {
		slog.Warn(errEnv.Error())
	}

	dsn := os.Getenv("DATABASE_URL")

	db, errDb := sql.Open("postgres", dsn)
	if errDb != nil {
		panic(errDb)
	}

	defer func(db *sql.DB) {
		errDbClose := db.Close()
		if errDbClose != nil {
			panic(errDbClose)
		}
	}(db)

	errPing := db.Ping()
	if errPing != nil {
		panic(errPing)
	}

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	defer func(rdb *redis.Client) {
		errRdb := rdb.Close()
		if errRdb != nil {
			panic(errRdb)
		}
	}(rdb)

	errRdbPing := rdb.Ping(ctx).Err()
	if errRdbPing != nil {
		panic(errRdbPing)
	}

	mux := http.NewServeMux()
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(os.Getenv("JWT_SECRET"), userRepo)
	userHandler := handlers.NewUserHandler(os.Getenv("JWT_SECRET"), userSvc)

	productRepo := repositories.NewProductRepository(db)
	productSvc := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	cartRepo := repositories.NewCartRepository(rdb)
	cartSvc := services.NewCartService(cartRepo)
	cartHandler := handlers.NewCartHandler(cartSvc)

	orderRepo := repositories.NewOrderRepository(db, rdb)
	orderSvc := services.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(&waitGroup, orderSvc)

	mux.HandleFunc("POST /api/register/", userHandler.RegisterUser)
	mux.HandleFunc("POST /api/login/", userHandler.LoginUser)

	mux.HandleFunc("GET /api/products/", productHandler.GetAllProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetProductByID)

	mux.Handle("POST /api/cart/add/", userHandler.AuthMiddleware(cartHandler.AddToCart))
	mux.Handle("POST /api/cart/remove/", userHandler.AuthMiddleware(cartHandler.RemoveFromCart))
	mux.Handle("GET /api/cart/", userHandler.AuthMiddleware(cartHandler.GetCart))

	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: middleware.CORSMiddleware(middleware.RateLimitMiddleware(mux)),
	}

	slog.Info("Server started", "port", ":"+os.Getenv("SERVER_PORT"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {

		errHTTP := srv.ListenAndServe()
		if errHTTP != nil {
			return
		}
	}()

	<-quit
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errShutdown := srv.Shutdown(ctx)
	if errShutdown != nil {
		panic(errShutdown)
	}

	waitGroup.Wait()
}
