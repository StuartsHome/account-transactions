## Overview

The server and database are containerised. The details of the two can be found in the file `docker-compose.yaml` located in the root of the project.

It is prefered to run the server locally and the database using a container. To run both in a container the file `store/store.go` needs to be modified on line 51 to set the variable `localServer` to `false`.

The database runs an entrypoint SQL file on startup to create the user, tables and some test data.  
The SQL file is located in `sql/init.sql`

The port (`8080`) for the server is located in `main.go`.  
The port (`3306`) for the database is located in the file `store/store.go`.  
*When running the DB in a container, the port in this file is overwritten with the port from `docker-compose.yaml`*

The project uses a MySQL docker image for the database. The details are as follows:
```sh
MYSQL_USER: storeuser
MYSQL_PASSWORD: example
MYSQL_ROOT_PASSWORD: example
MYSQL_HOST: "host.docker.internal"
MYSQL_DATABASE: store
```

## Getting Started
### Running the DB in a container and the server as a local binary:

>To start the DB in a container.  
1. `make start-db`

>To start the server as a binary file.  
2. `make start-local`

> To run tests. 
3. `make test`

> To view the documentation.
4. Access the link: `http://localhost:8080/swagger/index.html`

> To finish.
5. `make stop`

## Examples queries

> Create a new account with document number `123`
```sh
curl -XPOST "http://localhost:8080/accounts" \
-H "Content-Type: application/json" \
-d '{"document_number": "123"}'
```

> Get account with account ID `1`
```sh
curl -XGET "http://0.0.0.0:8080/accounts/1"
```

> Create a new payment transaction of `123.45`
```sh
curl -XPOST "http://0.0.0.0:8080/transactions" \
-H "Content-Type: application/json" \
-d '{"account_id": "1", "operation_type_id": 4, "amount": 123.45}'
```