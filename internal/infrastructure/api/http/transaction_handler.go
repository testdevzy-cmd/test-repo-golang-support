package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/test-repo-golang-support/internal/core/application/commands"
	"github.com/test-repo-golang-support/internal/core/application/services"
	"github.com/test-repo-golang-support/internal/core/domain/valueobjects"
)

// TransactionHandler handles HTTP requests for transactions
// Most complex multi-layer relationship:
// TransactionHandler -> DepositHandler/TransferHandler -> TransactionService ->
//   (TransactionRepository + AccountRepository) -> (Transaction + Account entities)
type TransactionHandler struct {
	depositHandler  *commands.DepositHandler
	transferHandler *commands.TransferHandler
	txService       *services.TransactionService
	logger          *log.Logger
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(
	depositHandler *commands.DepositHandler,
	transferHandler *commands.TransferHandler,
	txService *services.TransactionService,
	logger *log.Logger,
) *TransactionHandler {
	return &TransactionHandler{
		depositHandler:  depositHandler,
		transferHandler: transferHandler,
		txService:       txService,
		logger:          logger,
	}
}

// GetTransaction handles GET /transactions/{id}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txID := vars["id"]

	tx, err := h.txService.GetTransaction(r.Context(), txID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	// BUG: Accessing tx.AccountID directly
	// Knowledge graph should track: TransactionHandler -> Transaction.AccountID
	h.logger.Printf("Retrieved transaction: %s for account %s", tx.ID, tx.AccountID)

	h.respondJSON(w, http.StatusOK, tx)
}

// GetAccountTransactions handles GET /accounts/{id}/transactions
// Knowledge graph path:
// TransactionHandler -> TransactionService.GetAccountTransactions -> TransactionRepository.FindByAccountID
func (h *TransactionHandler) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	transactions, err := h.txService.GetAccountTransactions(r.Context(), accountID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, transactions)
}

// CreateDeposit handles POST /accounts/{id}/deposit
// Complex multi-layer command execution:
// HTTP -> DepositCommand -> DepositHandler -> TransactionService.CreateDeposit ->
//   (TransactionRepository + AccountRepository) -> (Transaction + Account)
func (h *TransactionHandler) CreateDeposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	var input struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create deposit command
	// Knowledge graph: TransactionHandler -> commands.DepositCommand -> valueobjects.Currency
	cmd := commands.DepositCommand{
		AccountID: accountID,
		Amount:    input.Amount,
		Currency:  valueobjects.Currency(input.Currency),
	}

	// Execute through command handler - multi-layer call
	tx, err := h.depositHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// BUG: Checking tx.IsCompleted() - if method is renamed, breaks here
	// Knowledge graph should detect: TransactionHandler -> Transaction.IsCompleted()
	if tx.IsCompleted() {
		h.logger.Printf("Deposit completed: %s, amount: %.2f", tx.ID, tx.Amount)
	}

	h.respondJSON(w, http.StatusCreated, tx)
}

// CreateTransfer handles POST /transfers
// Most complex relationship - involves 2 accounts and 1 transaction
func (h *TransactionHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SourceAccountID string  `json:"source_account_id"`
		TargetAccountID string  `json:"target_account_id"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create transfer command
	cmd := commands.TransferCommand{
		SourceAccountID: input.SourceAccountID,
		TargetAccountID: input.TargetAccountID,
		Amount:          input.Amount,
		Currency:        valueobjects.Currency(input.Currency),
	}

	// Execute through command handler
	// This triggers the most complex multi-layer relationship in the codebase:
	// TransactionHandler -> TransferHandler -> TransactionService.CreateTransfer ->
	//   AccountRepository.FindByID (2x) -> Account.IsActive (2x) -> Account.GetBalance (2x) ->
	//   Account.UpdateBalance (2x) -> TransactionRepository.Save -> Transaction.Complete
	tx, err := h.transferHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Access transaction fields - knowledge graph should track all field accesses
	h.logger.Printf("Transfer completed: %s, from %s to %s, amount: %.2f",
		tx.ID, tx.SourceAccountID, tx.TargetAccountID, tx.Amount)

	h.respondJSON(w, http.StatusCreated, tx)
}

// ProcessPendingTransactions handles POST /transactions/process
func (h *TransactionHandler) ProcessPendingTransactions(w http.ResponseWriter, r *http.Request) {
	if err := h.txService.ProcessPendingTransactions(r.Context()); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Pending transactions processed",
	})
}

// Helper methods
func (h *TransactionHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *TransactionHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// SetupTransactionRoutes configures transaction routes
func SetupTransactionRoutes(router *mux.Router, handler *TransactionHandler) {
	router.HandleFunc("/transactions/{id}", handler.GetTransaction).Methods("GET")
	router.HandleFunc("/transactions/process", handler.ProcessPendingTransactions).Methods("POST")
	router.HandleFunc("/accounts/{id}/transactions", handler.GetAccountTransactions).Methods("GET")
	router.HandleFunc("/accounts/{id}/deposit", handler.CreateDeposit).Methods("POST")
	router.HandleFunc("/transfers", handler.CreateTransfer).Methods("POST")
}

