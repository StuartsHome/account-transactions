package model

import "time"

type AccountImpl struct {
	AccountID      *int   `json:"account_id" db:"Account_ID"`
	DocumentNumber string `json:"document_number" db:"Document_Number"`
}

type OperationImpl struct {
	OperationTypeID int    `json:"operation_type_id" db:"OperationType_ID"`
	Description     string `json:"description" db:"Description"`
}

type Transactions []TransactionImpl
type TransactionImpl struct {
	TransactionID   *int       `json:"transaction_id" db:"Transaction_ID"`
	AccountID       int        `json:"account_id" db:"Account_ID"`
	OperationTypeID int        `json:"operation_type_id" db:"OperationType_ID"`
	Amount          float32    `json:"amount" db:"Amount"`
	Balance         float32    `json:"balance" db:"Balance"`
	EventDate       *time.Time `json:"-" db:"EventDate"`
}

func NewAccount(accountId *int, documentNumber string) *AccountImpl {
	return &AccountImpl{
		AccountID:      accountId,
		DocumentNumber: documentNumber,
	}
}

func NewTransaction(transactionId *int, accountId int, operationTypeId int, amount float32, balance float32, eventDate *time.Time) *TransactionImpl {
	return &TransactionImpl{
		TransactionID:   transactionId,
		AccountID:       accountId,
		OperationTypeID: operationTypeId,
		Amount:          amount,
		Balance:         balance,
		EventDate:       eventDate,
	}
}

func (t *OperationImpl) IsPurchase() bool {
	return t.OperationTypeID == 1
}

func (t *OperationImpl) IsPayment() bool {
	return t.OperationTypeID == 4
}

func ProcessNegativePayments(transactions Transactions, amount float32) (Transactions, float32, error) {
	currAmount := amount
	for i, transaction := range transactions {
		// If payment + amount > 0
		if currAmount > 0 && currAmount+transaction.Balance >= 0 {
			currAmount += transaction.Balance
			transactions[i].Balance = 0
		} else if currAmount > 0 && currAmount+transaction.Balance < 0 {
			// Partial payment
			transactions[i].Balance += currAmount
			currAmount = 0
			break
		}
	}

	return transactions, currAmount, nil
}
