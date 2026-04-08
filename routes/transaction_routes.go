package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"project-keuangan-keluarga/controller"
	"project-keuangan-keluarga/middleware"
)

func KeuanganRoutes(transactionsCtrl *controller.ControllerHandlerTransaction) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)         // custom structured JSON logger
	r.Use(middleware.MiddlewareAuth) // use
	r.Use(chimw.Recoverer)           // recover from panics
	r.Use(chimw.RequestID)           // inject X-Request-Id header

	// Health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// API v1 routes
	r.Post("/", transactionsCtrl.CreateNewTransactions_Bp)
	r.Put("/update", transactionsCtrl.UpdateTransactions_Bp)
	r.Get("/{id}", transactionsCtrl.GetTransactionById_Bp)
	r.Get("/all", transactionsCtrl.GetAllTransaction_Bp)
	r.Get("/avg-income-day", transactionsCtrl.GetAvgIncomeDay_Bp)
	r.Get("/avg-expense-day", transactionsCtrl.GetAvgExpenseDay_Bp)
	r.Get("/avg-income-week", transactionsCtrl.GetAvgIncomeWeek_Bp)
	r.Get("/avg-expense-week", transactionsCtrl.GetAvgExpenseWeek_Bp)
	r.Get("/avg-income-month", transactionsCtrl.GetAvgIncomeMonth_Bp)
	r.Get("/avg-expense-month", transactionsCtrl.GetAvgExpenseMonth_Bp)
	r.Get("/expense", transactionsCtrl.GetTransactionDataInExpenseType_Bp)
	r.Get("/income", transactionsCtrl.GetTransactionDataInIncomeType_Bp)
	r.Get("/expense-day-category", transactionsCtrl.GetAvgExpenseDayNameCategory_Bp)
	r.Get("/income-day-category", transactionsCtrl.GetAvgIncomeDayNameCategory_Bp)

	return r
}
