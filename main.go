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
	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/middleware/ratelimiter"
	"project-keuangan-keluarga/repository"
	"project-keuangan-keluarga/routes"
	"project-keuangan-keluarga/service"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// ── 1. Database ─────────────────────────────────────────────
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[MAIN] Failed to load .env file: %v", err)
	}

	dbCfg := config.DefaultDatabaseConfig()
	db, err := config.InitDB(dbCfg)
	if err != nil {
		log.Fatalf("[MAIN] Failed to initialize database: %v", err)
	}
	defer db.Close()

	// ── 2. Rate Limiter ─────────────────────────────────────────
	// In-memory store with cleanup every 2 minutes
	store := ratelimiter.NewMemoryStore(2 * time.Minute)
	defer store.Stop()

	// General limiter — 60 req/min (configurable via env)
	generalCfg := ratelimiter.LoadFromEnv()
	generalLimiter := ratelimiter.NewLimiter(generalCfg, store)
	log.Printf("[RATE-LIMIT] General: %d req/%s, burst=%d",
		generalCfg.Rate, generalCfg.Window, generalCfg.BurstCapacity)

	// Strict limiter — 10 req/min for auth endpoints (configurable via env)
	strictCfg := ratelimiter.LoadStrictFromEnv()
	strictLimiter := ratelimiter.NewLimiter(strictCfg, store)
	log.Printf("[RATE-LIMIT] Strict: %d req/%s, burst=%d",
		strictCfg.Rate, strictCfg.Window, strictCfg.BurstCapacity)

	// ── 3. Dependency Injection ──────────────────────────────────
	// ── 3. Repositories ──────────────────────────────────────────
	userRepo := repository.NewExampleRepository(db)
	familieRepo := repository.NewFamilieRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	familyMemberRepo := repository.NewFamilyMemberRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	goalsRepo := repository.NewGoalsRepository(db)

	// ── 4. Services ─────────────────────────────────────────────
	userSvc := service.NewUserService(userRepo)
	familieSvc := service.NewFamilieService(familieRepo, userRepo, familyMemberRepo, categoryRepo)
	familyMemberSvc := service.NewFamilyMemberService(familyMemberRepo, userRepo)
	categorySvc := service.NewCategoryService(categoryRepo, userRepo, familyMemberRepo)
	transactionSvc := service.NewTransactionService(transactionRepo, budgetRepo, familyMemberRepo)
	budgetSvc := service.NewBudgetService(budgetRepo, userRepo, familyMemberRepo)
	goalsSvc := service.NewGoalsService(goalsRepo, userRepo, familyMemberRepo)

	// ── 5. Controllers ──────────────────────────────────────────
	userCtrl := controller.NewUserController(userSvc)
	familieCtrl := controller.NewControllerHandlerFamilie(familieSvc)
	familyMemberCtrl := controller.NewControllerHandlerFamilyMember(familyMemberSvc)
	categoryCtrl := controller.NewControllerHandlerCategory(categorySvc)
	transactionCtrl := controller.NewControllerHandlerTransaction(transactionSvc)
	budgetCtrl := controller.NewBudgetController(budgetSvc)
	goalsCtrl := controller.NewControllerGoals(goalsSvc)

	// ── 6. Routes ────────────────────────────────────────────────
	route := chi.NewRouter()

	// CORS must be the very first middleware so it handles OPTIONS
	// preflight requests before auth or rate-limit middleware can
	// reject them.
	route.Use(middleware.CorsMiddleware())

	router := routes.UserRoutes(userCtrl, generalLimiter, strictLimiter)
	router_category := routes.CategoryRoutes(categoryCtrl, generalLimiter)
	router_transaction := routes.KeuanganRoutes(transactionCtrl, generalLimiter)
	router_budget := routes.BudgetRoutes(budgetCtrl, generalLimiter)
	router_goals := routes.GoalsRoutes(goalsCtrl, generalLimiter)
	router_familie := routes.FamilieRoutes(familieCtrl, generalLimiter)
	router_family_member := routes.FamilyMemberRoutes(familyMemberCtrl, generalLimiter)

	route.Mount("/api/v1/users", router)
	route.Mount("/api/v1/categories", router_category)
	route.Mount("/api/v1/transactions", router_transaction)
	route.Mount("/api/v1/budgets", router_budget)
	route.Mount("/api/v1/goals", router_goals)
	route.Mount("/api/v1/familie", router_familie)
	route.Mount("/api/v1/family-members", router_family_member)

	// Serving static profile images
	route.Handle("/uploadsProfile/*", http.StripPrefix("/uploadsProfile/", http.FileServer(http.Dir("./uploadsProfile"))))

	// ── 5. HTTP Server ───────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      route,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── 6. Graceful Shutdown ─────────────────────────────────────
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
