package server

import (
	"account-transactions/model"
	"account-transactions/store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// HandleGetAccount retrieves an account.
//
//	@Summary		Retrieves an account by ID
//	@Description	Retrieve an account with the provided account ID.
//	@Tags			account
//	@Accept			json
//	@Produce		json
//
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Success		200	{object}	model.AccountImpl
//
//	@Router			/accounts/{accountId} [get]
func HandleGetAccount(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get account ID from URL params.
		accountId := chi.URLParam(r, "accountId")
		// Convert string to int.
		accountIdInt, err := strconv.Atoi(accountId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(fmt.Appendf(nil, "invalid account ID %s: %v", accountId, err))
			return
		}

		gotAccount, err := db.GetAccount(accountIdInt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "account not found with ID %s: %v", accountId, err))
			return
		}

		// Success.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(gotAccount)
	}

}

// HandleAccountPost creates a new account.
//
//	@Summary		Create a new account
//	@Description	Creates an account with the provided document number.
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Body			{object} model.AccountImpl
//
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Success		201	{object}	model.AccountImpl
//
//	@Router			/accounts [post]
func HandleAccountPost(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := model.AccountImpl{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "err %v", err))
			return
		}
		if err := json.Unmarshal([]byte(body), &account); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "err %v", err))
			return
		}
		newAccount, err := db.CreateAccount(account.DocumentNumber)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "err %v", err))
			return
		}

		// Success.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAccount)
	}
}

// HandleTransactionPost creates a new transaction.
//
//	@Summary		Create a new transaction
//	@Description	Creates a transaction with the provided account ID, operation type ID, and amount.
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Body			model.TransactionImpl	true				"Transaction to create"
//
//	@Failure		404						{string}	string	"Not Found"
//	@Failure		500						{string}	string	"Internal Server Error"
//	@Success		201						{object}	model.TransactionImpl
//
//	@Router			/transactions [post]
func HandleTransactionPost(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		transaction := model.TransactionImpl{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "err %v", err))
			return
		}
		if err := json.Unmarshal([]byte(body), &transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "err %v", err))
			return
		}

		// Validate account id.
		_, err = db.GetAccount(transaction.AccountID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(fmt.Appendf(nil, "err account doesn't exist %v", err))
			return
		}

		// Validate operation id.
		_, err = db.GetOperation(transaction.OperationTypeID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(fmt.Appendf(nil, "err operation doesn't exist %v", err))
			return
		}

		// Store.
		result, err := db.CreateTransaction(transaction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "err %v", err))
			return
		}

		// Success.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}
}
