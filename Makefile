# Include variables from the .envrc file
include .envrc
run:
	go run ./cmd/api -dsn=${VENDORS_DB_DSN}

psql:
	psql ${VENDORS_DB_DSN}

up:
	@echo 'Running up migrations...'
	migrate -path=migrations -database=${VENDORS_DB_DSN} up

down:
	@echo 'Running down migrations...'
	migrate -path=migrations -database=${VENDORS_DB_DSN} down

version:
	@echo 'Migrating to version 1'
	migrate -path=migrations -database=${VENDORS_DB_DSN} goto 0
