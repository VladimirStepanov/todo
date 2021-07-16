rebuild:
	@docker-compose -f docker/docker-compose-prod.yml up --build

prod:
	@docker-compose -f docker/docker-compose-prod.yml up -d

migrateup:
	@/bin/bash scripts/migrate.sh up

migratedown:
	@/bin/bash scripts/migrate.sh down

test:
	@docker-compose -f docker/docker-compose-test.yml up --build

.PHONY: rebuild prod migrateup migratedown test