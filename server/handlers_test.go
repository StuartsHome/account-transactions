package server

import (
	mock_store "account-transactions/mocks"
	"account-transactions/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	accountId      = "123"
	accountIdInt   = 123
	documentNumber = "20251027"
	transactionID  = 111
)

func TestHandleGetAccount(t *testing.T) {
	// Given.
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	// Create a chi Context object
	chiCtx := chi.NewRouteContext()

	// Create a new test request with the additional Chi contetx
	reqWithCtx := req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("accountId", accountId)

	recorder := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	m := mock_store.NewMockStore(ctrl)
	m.EXPECT().
		GetAccount(accountId).
		Return(&model.AccountImpl{
			AccountID:      model.IntToPtr(123),
			DocumentNumber: documentNumber,
		}, nil)

	// When.
	// This is the handler func we want to test
	hf := http.HandlerFunc(HandleGetAccount(m))
	hf.ServeHTTP(recorder, reqWithCtx)

	// Then.
	// Check the status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is correct
	expected := fmt.Sprintf("{\"account_id\":%d,\"document_number\":\"%s\"}\n", accountIdInt, documentNumber)
	got := recorder.Body.String()
	assert.Equal(t, expected, got)
}

func TestHandleAccountPost(t *testing.T) {
	// Given.
	body := fmt.Sprintf("{\"document_number\":\"%s\"}", documentNumber)
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	m := mock_store.NewMockStore(ctrl)
	m.EXPECT().
		CreateAccount(documentNumber).
		Return(&model.AccountImpl{
			AccountID:      model.IntToPtr(accountIdInt),
			DocumentNumber: documentNumber,
		}, nil)

	// When.
	// This is the handler func we want to test
	hf := http.HandlerFunc(HandleAccountPost(m))
	hf.ServeHTTP(recorder, req)

	// Then.
	// Check the status code
	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is correct
	expected := fmt.Sprintf("{\"account_id\":%d,\"document_number\":\"%s\"}\n", accountIdInt, documentNumber)
	got := recorder.Body.String()
	assert.Equal(t, expected, got)
}

func TestHandleTransactionPost(t *testing.T) {
	// Given.
	transaction := model.NewTransaction(&transactionID, accountId, 4, 5000.00, nil)

	marshalledTransaction, err := json.Marshal(transaction)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/", strings.NewReader(string(marshalledTransaction)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	m := mock_store.NewMockStore(ctrl)
	m.EXPECT().
		GetAccount(accountId).
		Return(&model.AccountImpl{
			AccountID:      &accountIdInt,
			DocumentNumber: documentNumber,
		}, nil)
	m.EXPECT().
		GetOperation(4).
		Return(&model.OperationsTypes{
			OperationTypeID: 4,
			Description:     "PAYMENT",
		}, nil)
	m.EXPECT().
		CreateTransaction(*transaction).
		Return(&model.TransactionImpl{
			TransactionID:   &transactionID,
			AccountID:       accountId,
			OperationTypeID: 4,
			Amount:          5000.00,
		}, nil)

	// When.
	// This is the handler func we want to test
	hf := http.HandlerFunc(HandleTransactionPost(m))
	hf.ServeHTTP(recorder, req)

	// Then.
	// Check the status code
	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is correct
	expected := string(marshalledTransaction) + "\n"
	got := recorder.Body.String()
	assert.Equal(t, expected, got)
}
