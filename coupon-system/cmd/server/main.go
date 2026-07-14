package main

import (
	"coupon-system/internal/api/handlers"
	"coupon-system/internal/config"
	"coupon-system/internal/repository"
	"coupon-system/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize repository (in-memory for this MVP)
	repo := repository.NewCouponRepository()

	// Wrap repository with cache
	cachedRepo := repository.NewCachedCouponRepository(repo, 5*time.Minute)

	// Initialize service
	couponService := service.NewCouponService(cachedRepo)

	// Set the service in handlers
	handlers.CouponSvc = couponService

	// Initialize router
	router := mux.NewRouter()

	// Register routes
	handlers.RegisterRoutes(router)

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("error writing response: %v", err)
		}
	}).Methods(http.MethodGet)

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Starting server on %s", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
