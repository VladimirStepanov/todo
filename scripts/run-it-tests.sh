#!/bin/bash

sh scripts/wait-postgres.sh
sh scripts/migrate.sh up
go test -v ./it