package main

import (
	"log"
	"net/http"

	"blog-api/internal/config"
	"blog-api/internal/database"
	"blog-api/internal/handlers"
	"blog-api/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load konfigurasi
	cfg := config.LoadConfig()

	// Koneksi ke database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Jalankan migration dan seed
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := database.SeedData(); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	// Setup router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Swagger documentation
	router.HandleFunc("/swagger.json", handlers.SwaggerJSON).Methods("GET")
	router.HandleFunc("/swagger", handlers.SwaggerUI).Methods("GET")

	// API routes (akan ditambahkan di langkah selanjutnya)
	api := router.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")

	// Post routes (protected)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	protected.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	protected.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")

	// Public post routes
	api.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", handlers.GetPost).Methods("GET")

	// Comment routes (protected)
	protected.HandleFunc("/posts/{post_id}/comments", handlers.CreateComment).Methods("POST")
	protected.HandleFunc("/posts/{post_id}/comments/{comment_id}", handlers.DeleteComment).Methods("DELETE")

	// Public comment routes
	api.HandleFunc("/posts/{post_id}/comments", handlers.GetComments).Methods("GET")

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
