package model

import "time"

type Account interface {
	GetAccount(int) (*AccountImpl, error)
	CreateAccount(string) (*AccountImpl, error)
}

type AccountImpl struct {
	AccountID      *int   `json:"account_id" db:"Account_ID"`
	DocumentNumber string `json:"document_number" db:"Document_Number"`
}

type Operation interface {
	GetOperation(int) (*OperationsTypes, error)
}

type OperationsTypes struct {
	OperationTypeID int    `json:"operation_type_id" db:"OperationType_ID"`
	Description     string `json:"description" db:"Description"`
}

type Transaction interface {
	GetTransaction(string) (*TransactionImpl, error)
	CreateTransaction(TransactionImpl) (*TransactionImpl, error)
}

type TransactionImpl struct {
	TransactionID   *int       `json:"transaction_id" db:"Transaction_ID"`
	AccountID       int        `json:"account_id" db:"Account_ID"`
	OperationTypeID int        `json:"operation_type_id" db:"OperationType_ID"`
	Amount          float32    `json:"amount" db:"Amount"`
	EventDate       *time.Time `json:"-" db:"EventDate"`
}

func NewAccount(accountId *int, documentNumber string) *AccountImpl {
	return &AccountImpl{
		AccountID:      accountId,
		DocumentNumber: documentNumber,
	}
}

func NewTransaction(transactionId *int, accountId int, operationTypeId int, amount float32, eventDate *time.Time) *TransactionImpl {
	return &TransactionImpl{
		TransactionID:   transactionId,
		AccountID:       accountId,
		OperationTypeID: operationTypeId,
		Amount:          amount,
		EventDate:       eventDate,
	}
}
