# egw-be
This project represents the EGW backend server.

## Setup
Install dependencies:
1. `go install`
Copy the .env vars:
2. `make setup`
Start docker dependencies:
3. `docker-compose up`
Migrate the DB:
4. `make migrate`
Run the server:
5. `make run`

The server is now running locally and listening for requests. 

## Testing
Ensure you have the Postgres database up, by running `docker-compose up`
Then run `make test` to execute all unit tests