#!/bin/bash

migrate -path migrations/ -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$MIGRATE_DB_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable $1