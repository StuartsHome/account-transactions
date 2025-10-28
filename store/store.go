package store

import (
	"account-transactions/model"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	GetAccount(int) (*model.AccountImpl, error)
	GetOperation(int) (*model.OperationsTypes, error)
	GetTransaction(string) (*model.TransactionImpl, error)
	CreateAccount(string) (*model.AccountImpl, error)
	CreateTransaction(model.TransactionImpl) (*model.TransactionImpl, error)
}

var _ Store = &StoreImpl{}

type StoreImpl struct {
	db *sqlx.DB
	Store
}

var dbport = 3306

func New() *StoreImpl {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("using store: mysql")
	return &StoreImpl{
		db: db,
	}
}

func connect() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	localServer := true
	if localServer {
		// When running the DB in a container and the server as a local binary.
		db, err = sqlx.Open("mysql", fmt.Sprintf("storeuser:example@tcp(0.0.0.0:%d)/store", dbport))
	} else {
		// When both the DB and server are running in containers.
		db, err = sqlx.Open("mysql", fmt.Sprintf("storeuser:example@tcp(db:%d)/store", dbport))
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("mysql connected on %d\n", dbport)
	return db, err
}
