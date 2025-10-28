package store

import (
	"account-transactions/model"
	"database/sql"
	"fmt"
)

func (s *StoreImpl) GetAccount(accountId string) (*model.AccountImpl, error) {

	var account model.AccountImpl
	err := s.db.Get(&account, "SELECT Account_ID, Document_Number FROM Accounts WHERE Account_ID=?", accountId)
	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("no account with id %s, err: %v", accountId, err)
	case err != nil:
		err = fmt.Errorf("query error: %v", err)
	}
	return &account, err
}

func (s *StoreImpl) GetOperation(operationId int) (*model.OperationsTypes, error) {

	var account model.OperationsTypes
	err := s.db.Get(&account, "SELECT OperationType_ID, Description FROM OperationsTypes WHERE OperationType_ID=?", operationId)
	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("no operation with id %d, err: %v", operationId, err)
	case err != nil:
		err = fmt.Errorf("query error: %v", err)
	}
	return &account, err
}

func (s *StoreImpl) GetTransaction(transactionId string) (*model.TransactionImpl, error) {

	var transaction model.TransactionImpl
	err := s.db.Get(&transaction, "SELECT Transaction_ID, Account_ID, OperationType_ID, Amount FROM Transactions WHERE Transaction_ID=?", transactionId)
	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("no transaction with id %s, err: %v", transactionId, err)
	case err != nil:
		err = fmt.Errorf("query error: %v", err)
	}
	return &transaction, err

}
func (s *StoreImpl) CreateAccount(docNumber string) (*model.AccountImpl, error) {

	stmt, err := s.db.Prepare("INSERT INTO Accounts(Document_Number) VALUES( ? )")
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	res, err := stmt.Exec(docNumber)
	if err != nil {
		return nil, err
	}
	// Get the transaction id from the inserted row.
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	lastIdConverted := int(lastId)
	account := model.NewAccount(&lastIdConverted, docNumber)
	return account, err
}

func (s *StoreImpl) CreateTransaction(transaction model.TransactionImpl) (*model.TransactionImpl, error) {
	stmt, err := s.db.Prepare("INSERT INTO Transactions(Account_ID, OperationType_ID, Amount, EventDate) VALUES( ?, ?, ?, Now() )")
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	res, err := stmt.Exec(transaction.AccountID, transaction.OperationTypeID, transaction.Amount)
	if err != nil {
		return nil, err
	}
	// Get the transaction id from the inserted row.
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	transactionId := int(lastId)
	resultTransaction := model.NewTransaction(&transactionId, transaction.AccountID, transaction.OperationTypeID, transaction.Amount, nil)

	return resultTransaction, err
}
