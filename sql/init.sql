DROP DATABASE IF EXISTS store;
CREATE DATABASE store;
CREATE USER 'storeuser'@'localhost' IDENTIFIED BY 'example';
GRANT ALL ON store.* TO 'storeuser'@'localhost';
USE store;

DROP TABLE IF EXISTS Accounts;
CREATE TABLE Accounts (
    Account_ID int NOT NULL auto_increment,
    Document_Number BIGINT (16) NOT NULL,
    PRIMARY KEY (Account_ID)
);
INSERT INTO Accounts ( Document_Number )
VALUES
(12345678900);

DROP TABLE IF EXISTS OperationsTypes;
CREATE TABLE OperationsTypes (
    OperationType_ID int NOT NULL auto_increment,
    Description VARCHAR (255) NOT NULL,
    PRIMARY KEY (OperationType_ID)
);
INSERT INTO OperationsTypes ( Description )
VALUES
("PURCHASE"),
("INSTALLMENT PURCHASE"),
("WITHDRAWAL"),
("PAYMENT");

DROP TABLE IF EXISTS Transactions;
CREATE TABLE Transactions (
    Transaction_ID int NOT NULL auto_increment,
    Account_ID int NOT NULL,
    OperationType_ID int NOT NULL,
    Amount DECIMAL (18,2) NOT NULL,
    Balance DECIMAL (18,2) DEFAULT 0.00 NOT NULL,
    EventDate DATETIME NOT NULL,
    PRIMARY KEY (Transaction_ID),
    FOREIGN KEY (Account_ID) REFERENCES Accounts(Account_ID),
    FOREIGN KEY (OperationType_ID) REFERENCES OperationsTypes(OperationType_ID)
);
