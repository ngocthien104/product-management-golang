package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"mini-product-management/backend/internal/handler"
)

func main() {
	_ = godotenv.Load()

	dbUser := env("DB_USER", "root")
	dbPassword := env("DB_PASSWORD", "")
	dbHost := env("DB_HOST", "127.0.0.1")
	dbPort := env("DB_PORT", "3306")
	dbName := env("DB_NAME", "mini_product_management")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("open MySQL: ", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to MySQL: %v; check .env and make sure MySQL is running", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	productHandler := handler.NewProductHandler(db)
	products := api.Group("/products")
	{
		products.GET("", productHandler.List)
		products.GET("/:slug", productHandler.GetBySlug)
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}
	api.GET("/categories", productHandler.Categories)

	port := env("PORT", "8080")
	log.Printf("API listening on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
