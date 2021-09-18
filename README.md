# Description

Simple REST API TODO application written using Clean Architecture principles

List of used libraries:

* [logrus](https://github.com/sirupsen/logrus) for logging
* [gin](https://github.com/gin-gonic/gin) web framework
* [cleanenv](https://github.com/ilyakaznacheev/cleanenv) for reading config
* [sqlx](https://github.com/jmoiron/sqlx) for working with DB
* [testify](https://github.com/stretchr/testify) for testing (mock, require, suite)
* [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) for database mocking
* [redis-mock](https://github.com/go-redis/redismock) for redis mocking
* [golang-jwt/jwt](https://github.com/golang-jwt/jwt) for JWT auth

## Startup configuration [file .env in root]

```env
#for confirm registration

EMAIL=confirmemail@gmail.com
EMAIL_PASSWORD=email_password

DOMAIN=site_domain

APP_ADDR=bind_addr
APP_PORT=bind_port

#db addr for migration
MIGRATE_DB_HOST=migrate_db_addr

#database
POSTGRES_HOST=todo_db
POSTGRES_PORT=5433
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=todo

#jwt keys
JWT_ACCESS_KEY=access_key
JWT_REFRESH_KEY=refresh_key

#redis
REDIS_HOST=tokendb
REDIS_PORT=6380

#other
MAX_LOGGED_IN=6
```

## Run

```bash
make rebuild
```

## Migrations

```bash
make migrateup
make migratedown
```

## Tests

```bash
make test #unit testing
make test.integrations #integration testing
```

## Documentation

Install  [swag](https://github.com/swaggo/swag/cmd/swag)

```bash
make swag #generate docs
```

See `/swagger/index.html` path.