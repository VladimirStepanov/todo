rebuild:
	@docker-compose -f docker/docker-compose-prod.yml up --build

prod:
	@docker-compose -f docker/docker-compose-prod.yml up -d

swag:
	@echo "package kk" > dummy.go
	@swag init --parseInternal=true -g cmd/main.go
	@rm -f dummy.go

migrateup:
	@/bin/bash scripts/migrate.sh up

migratedown:
	@/bin/bash scripts/migrate.sh down

test:
	@docker-compose -f docker/docker-compose-test.yml up --build

test.integrations:
	@docker-compose -f docker/docker-compose-it-tests.yml up --build --abort-on-container-exit --exit-code-from it_test_todo

.PHONY: rebuild prod migrateup migratedown test test.integrations