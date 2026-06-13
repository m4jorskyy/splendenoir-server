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
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	mux := http.NewServeMux()
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(os.Getenv("JWT_SECRET"), userRepo)
	userHandler := handlers.NewUserHandler(os.Getenv("JWT_SECRET"), userSvc)

	productRepo := repositories.NewProductRepository(db)
	productSvc := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	mux.HandleFunc("POST /api/register/", userHandler.RegisterUser)
	mux.HandleFunc("POST /api/login/", userHandler.LoginUser)

	mux.HandleFunc("GET /api/products/", productHandler.GetAllProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetProductByID)

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
}
