#!/bin/bash

sh scripts/wait-postgres.sh
sh scripts/migrate.sh "drop -f"
sh scripts/migrate.sh up
go test -v ./it