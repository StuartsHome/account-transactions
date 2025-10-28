package store

import (
	"account-transactions/model"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accountId        = "123"
	accountIdInt     = 123
	invalidAccountId = 999
	documentNumber   = "20251027"
	transactionID    = 111
)

func TestGetAccount_Success(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"Account_ID", "Document_Number"}).
		AddRow(accountIdInt, documentNumber)

	mock.ExpectQuery("SELECT Account_ID, Document_Number FROM Accounts WHERE Account_ID=?").
		WithArgs(accountIdInt).
		WillReturnRows(rows)

	// When.
	account, err := store.GetAccount(accountIdInt)

	// Then.
	require.NoError(t, err)
	expectedAccount := &model.AccountImpl{
		AccountID:      model.IntToPtr(accountIdInt),
		DocumentNumber: documentNumber,
	}
	assert.Equal(t, expectedAccount, account)
}

func TestGetAccount_NotFound(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	mock.ExpectQuery("SELECT Account_ID, Document_Number FROM Accounts WHERE Account_ID=?").
		WithArgs(invalidAccountId).
		WillReturnError(sql.ErrNoRows)

	// When.
	account, err := store.GetAccount(invalidAccountId)

	// Then.
	require.Error(t, err)
	assert.Equal(t, model.NewAccount(nil, ""), account)
	assert.Contains(t, err.Error(), fmt.Sprintf("no account with id %d", invalidAccountId))
}

func TestCreateAccount_Success(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	mock.ExpectPrepare("INSERT INTO Accounts\\(Document_Number\\) VALUES\\( \\? \\)").
		ExpectExec().
		WithArgs("20251027").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// When.
	account, err := store.CreateAccount("20251027")

	// Then.
	require.NoError(t, err)
	expectedAccount := &model.AccountImpl{
		AccountID:      model.IntToPtr(1),
		DocumentNumber: "20251027",
	}
	assert.Equal(t, expectedAccount, account)
}

func TestCreateAccount_Fail(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	mock.ExpectPrepare("INSERT INTO Accounts\\(Document_Number\\) VALUES\\( \\? \\)").
		ExpectExec().
		WithArgs("20251027").
		WillReturnError(sql.ErrConnDone)

	// When.
	account, err := store.CreateAccount("20251027")

	// Then.
	require.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "sql: connection is already closed")
}

func TestCreateTransaction_Success(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Transactions(Account_ID, OperationType_ID, Amount, EventDate) VALUES( ?, ?, ?, Now() )`)).
		ExpectExec().
		WithArgs(accountIdInt, 4, 5000.00).
		WillReturnResult(sqlmock.NewResult(int64(transactionID), 1))

	// When.
	transaction, err := store.CreateTransaction(*model.NewTransaction(nil, accountIdInt, 4, 5000.00, nil))

	// Then.
	require.NoError(t, err)
	expectedTransaction := &model.TransactionImpl{
		TransactionID:   &transactionID,
		AccountID:       accountIdInt,
		OperationTypeID: 4,
		Amount:          5000.00,
	}
	assert.Equal(t, expectedTransaction, transaction)
}

func TestCreateTransaction_Fail(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Transactions(Account_ID, OperationType_ID, Amount, EventDate) VALUES( ?, ?, ?, Now() )`)).
		ExpectExec().
		WithArgs(accountIdInt, 4, 5000.00).
		WillReturnError(sql.ErrConnDone)

	// When.
	transaction, err := store.CreateTransaction(
		*model.NewTransaction(nil, accountIdInt, 4, 5000.00, nil),
	)

	// Then.
	require.Error(t, err)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "sql: connection is already closed")
}
