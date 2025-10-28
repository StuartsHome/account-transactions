package store

import (
	"account-transactions/model"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccount_Success(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := &StoreImpl{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"Account_ID", "Document_Number"}).
		AddRow(123, "20251027")

	mock.ExpectQuery("SELECT Account_ID, Document_Number FROM Accounts WHERE Account_ID=?").
		WithArgs("123").
		WillReturnRows(rows)

	// When.
	account, err := store.GetAccount("123")

	// Then.
	require.NoError(t, err)
	expectedAccount := &model.AccountImpl{
		AccountID:      model.IntToPtr(123),
		DocumentNumber: "20251027",
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
		WithArgs("999").
		WillReturnError(sql.ErrNoRows)

	// When.
	account, err := store.GetAccount("999")

	// Then.
	require.Error(t, err)
	assert.Equal(t, model.NewAccount(nil, ""), account)
	assert.Contains(t, err.Error(), "no account with id 999")
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
		WithArgs("123", 4, 5000.00).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// When.
	transaction, err := store.PutTransaction(*model.NewTransaction(nil, "123", 4, 5000.00, nil))

	// Then.
	require.NoError(t, err)
	expectedTransaction := &model.TransactionImpl{
		TransactionID:   model.IntToPtr(1),
		AccountID:       "123",
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
		WithArgs("123", 4, 5000.00).
		WillReturnError(sql.ErrConnDone)

	// When.
	transaction, err := store.PutTransaction(*model.NewTransaction(nil, "123", 4, 5000.00, nil))

	// Then.
	require.Error(t, err)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "sql: connection is already closed")
}
