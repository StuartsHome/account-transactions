package server

import (
	mock_store "account-transactions/mocks"
	"account-transactions/model"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandleGetAccount(t *testing.T) {
	// Given.
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a chi Context object
	chiCtx := chi.NewRouteContext()

	// Create a new test request with the additional Chi contetx
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("accountId", "123")

	recorder := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	m := mock_store.NewMockStore(ctrl)
	m.EXPECT().GetAccount("123").Return(&model.AccountImpl{
		AccountID:      model.IntToPtr(123),
		DocumentNumber: "20251027",
	}, nil)

	// When.
	// This is the handler func we want to test
	hf := http.HandlerFunc(HandleGetAccount(m))
	hf.ServeHTTP(recorder, req)

	// Then.
	// Check the status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is correct
	expected := "{\"account_id\":123,\"document_number\":\"20251027\"}\n"
	actual := recorder.Body.String()
	assert.Equal(t, expected, actual)
}

func TestHandleAccountPost(t *testing.T) {
	// Given.
	body := `{"document_number":"123"}`
	r, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	m := mock_store.NewMockStore(ctrl)
	m.EXPECT().CreateAccount("123").Return(&model.AccountImpl{
		AccountID:      model.IntToPtr(123),
		DocumentNumber: "20251027",
	}, nil)

	// When.
	// This is the handler func we want to test
	hf := http.HandlerFunc(HandleAccountPost(m))
	hf.ServeHTTP(recorder, r)

	// Then.
	// Check the status code
	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is correct
	expected := "{\"account_id\":123,\"document_number\":\"20251027\"}\n"
	actual := recorder.Body.String()
	assert.Equal(t, expected, actual)
}

func TestHandleTransactionPost(t *testing.T) {
	// Given.
	transaction := model.NewTransaction(model.IntToPtr(123), "123", 4, 5000.00, nil)

	marshalledTransaction, _ := json.Marshal(transaction)
	r, err := http.NewRequest("POST", "/", strings.NewReader(string(marshalledTransaction)))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	m := mock_store.NewMockStore(ctrl)
	m.EXPECT().
		GetAccount("123").Return(&model.AccountImpl{
		AccountID:      model.IntToPtr(123),
		DocumentNumber: "20251027",
	}, nil)
	m.EXPECT().GetOperation(4).Return(&model.OperationsTypes{
		OperationTypeID: 4,
		Description:     "PAYMENT",
	}, nil)
	m.EXPECT().
		CreateTransaction(*transaction).Return(&model.TransactionImpl{
		TransactionID:   model.IntToPtr(123),
		AccountID:       "123",
		OperationTypeID: 4,
		Amount:          5000.00,
	}, nil)

	// When.
	// This is the handler func we want to test
	hf := http.HandlerFunc(HandleTransactionPost(m))
	hf.ServeHTTP(recorder, r)

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
