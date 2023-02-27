# Include variables from the .envrc file
include .envrc
run:
	go run ./cmd/api -dsn=${VENDORS_DB_DSN}

psql:
	psql ${VENDORS_DB_DSN}

up:
	migrate -path=migrations -database=$VENDORS_DB_DSN up

down:
	migrate -path=migrations -database=$VENDORS_DB_DSN down
