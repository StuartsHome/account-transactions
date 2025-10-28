package model

import "time"

type AccountImpl struct {
	AccountID      *int   `json:"account_id" db:"Account_ID"`
	DocumentNumber string `json:"document_number" db:"Document_Number"`
}

type OperationsTypes struct {
	OperationTypeID int    `json:"OperationType_ID" db:"OperationType_ID"`
	Description     string `json:"Description" db:"Description"`
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
