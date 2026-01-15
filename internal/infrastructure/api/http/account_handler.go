package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/test-repo-golang-support/internal/core/application/commands"
	"github.com/test-repo-golang-support/internal/core/application/services"
	"github.com/test-repo-golang-support/internal/core/domain/entities"
	"github.com/test-repo-golang-support/internal/core/domain/valueobjects"
)

// AccountHandler handles HTTP requests for accounts
// Knowledge graph multi-layer relationship:
// AccountHandler -> commands.CreateAccountHandler -> services.AccountService -> 
//   repositories.AccountRepository -> entities.Account
type AccountHandler struct {
	createHandler  *commands.CreateAccountHandler
	suspendHandler *commands.SuspendAccountHandler
	accountService *services.AccountService
	logger         *log.Logger
}

// NewAccountHandler creates a new AccountHandler
// Knowledge graph should track constructor dependencies
func NewAccountHandler(
	createHandler *commands.CreateAccountHandler,
	suspendHandler *commands.SuspendAccountHandler,
	accountService *services.AccountService,
	logger *log.Logger,
) *AccountHandler {
	return &AccountHandler{
		createHandler:  createHandler,
		suspendHandler: suspendHandler,
		accountService: accountService,
		logger:         logger,
	}
}

// GetAccount handles GET /accounts/{id}
// 5-layer deep relationship:
// HTTP Handler -> AccountService -> AccountRepository -> Account entity
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	account, err := h.accountService.GetAccount(r.Context(), accountID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	// BUG: Accessing account.Email directly
	// If entities.Account.Email is renamed, knowledge graph should detect
	// the break in this 5-layer chain
	h.logger.Printf("Retrieved account: %s (%s)", account.ID, account.Email)

	h.respondJSON(w, http.StatusOK, account)
}

// CreateAccount handles POST /accounts
// Multi-layer command pattern:
// HTTP -> CreateAccountCommand -> CreateAccountHandler -> AccountService -> Repository -> Entity
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OwnerID     string `json:"owner_id"`
		Email       string `json:"email"`        // JSON field name
		AccountType string `json:"account_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create command - knowledge graph tracks command structure
	cmd := commands.CreateAccountCommand{
		OwnerID:     input.OwnerID,
		Email:       input.Email,
		AccountType: entities.AccountType(input.AccountType),
	}

	// Execute through command handler - 4-layer deep call
	account, err := h.createHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, account)
}

// GetAccountBalance handles GET /accounts/{id}/balance
// Knowledge graph path:
// AccountHandler -> AccountService.GetAccountBalance -> AccountRepository.GetBalance -> Account.GetBalance()
// Also involves: valueobjects.Money
func (h *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	balance, err := h.accountService.GetAccountBalance(r.Context(), accountID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	// BUG: Accessing Money.Amount directly vs using String() method
	// Knowledge graph should track: AccountHandler -> valueobjects.Money.Amount
	response := map[string]interface{}{
		"account_id": accountID,
		"balance":    balance.Amount,
		"currency":   balance.Currency,
		"formatted":  balance.String(),
	}

	h.respondJSON(w, http.StatusOK, response)
}

// SuspendAccount handles POST /accounts/{id}/suspend
// Knowledge graph path:
// AccountHandler -> SuspendAccountHandler -> AccountService.SuspendAccount -> Account.Suspend()
func (h *AccountHandler) SuspendAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	cmd := commands.SuspendAccountCommand{
		AccountID: accountID,
		Reason:    "Manual suspension via API",
	}

	if err := h.suspendHandler.Handle(r.Context(), cmd); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Account suspended successfully",
	})
}

// Helper methods
func (h *AccountHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AccountHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// SetupAccountRoutes configures account routes
func SetupAccountRoutes(router *mux.Router, handler *AccountHandler) {
	router.HandleFunc("/accounts", handler.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", handler.GetAccount).Methods("GET")
	router.HandleFunc("/accounts/{id}/balance", handler.GetAccountBalance).Methods("GET")
	router.HandleFunc("/accounts/{id}/suspend", handler.SuspendAccount).Methods("POST")
}

