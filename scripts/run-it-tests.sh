#!/bin/bash

sh scripts/wait-postgres.sh
go run cmd/migrate/migrate.go -direction down
go run cmd/migrate/migrate.go -direction up
go test -v ./it