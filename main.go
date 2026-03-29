package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-keuangan-keluarga/config"
	"project-keuangan-keluarga/controller"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/routes"
	"project-keuangan-keluarga/service"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// ── 1. Database ──────────────────────────────────────────────

	if err := godotenv.Load(); err != nil {
		log.Fatalf("[MAIN] Failed to load .env file: %v", err)
	}

	dbCfg := config.DefaultDatabaseConfig()
	db, err := config.InitDB(dbCfg)
	if err != nil {
		log.Fatalf("[MAIN] Failed to initialize database: %v", err)
	}
	defer db.Close()

	// ── 2. Dependency Injection ──────────────────────────────────
	userRepo := repository.NewExampleRepository(db)
	userSvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSvc)

	// keuangan injection
	keuanganRepo := repository.NewKeuanganRepository(db)
	keuanganSvc := service.NewKeuanganService(keuanganRepo)
	keuanganCtrl := controller.NewControllerHandlerKeuangan(keuanganSvc)

	// ── 3. Routes ────────────────────────────────────────────────
	route := chi.NewRouter()
	subRoiter := route.With()
	router := routes.UserRoutes(userCtrl)
	router_keuangan := routes.KeuanganRoutes(keuanganCtrl)

	subRoiter.Mount("/api/v1/users", router)
	subRoiter.Mount("/api/v1/keuangans", router_keuangan)

	// ── 4. HTTP Server ───────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      route,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── 5. Graceful Shutdown ─────────────────────────────────────
	go func() {
		log.Println("[SERVER] Listening on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[SERVER] Failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[SERVER] Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[SERVER] Forced shutdown: %v", err)
	}

	log.Println("[SERVER] Server stopped")
}
