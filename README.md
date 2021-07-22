# Description

Simple REST API TODO application written using Clean Architecture principles

List of used libraries:

* [logrus](github.com/sirupsen/logrus) for logging
* [gin](github.com/gin-gonic/gin) web framework
* [cleanenv](github.com/ilyakaznacheev/cleanenv) for reading config
* [sqlx](github.com/jmoiron/sqlx) for working with DB
* [testify](github.com/stretchr/testify) for testing (mock, require)
* [go-sqlmock](github.com/DATA-DOG/go-sqlmock) for database mocking

## Startup configuration [file .env in root]

```env
#for confirm registration

EMAIL=confirmemail@gmail.com
EMAIL_PASSWORD=email_password

DOMAIN=site_domain

APP_ADDR=bind_addr
APP_PORT=bind_port

MIGRATE_DB_HOST=migrate_db_addr #db addr for migration

#database
POSTGRES_HOST=todo_db
POSTGRES_PORT=5433
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=todo
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

# Endpoints

## Registration

`POST /auth/sign-up`

Params (json):
* email [string]
* password [string]

#### Request
```bash
curl -L -X POST 'localhost:8080/auth/sign-up' -H 'Content-Type: application/json' --data-raw '{
    "email": "test1234test@mmail.com",
    "password": "1234567891"
}'
```

#### Response

```json
{
    "status":"success",
}
```

## Email confirmation

`GET /auth/confirm/:activated_link`

Params (URL):
* activated_link [string]

#### Request
```bash
curl http://localhost:8080/auth/confirm/238930b6-b3f3-461c-9cbd-f59e7e6bf072
```

#### Response

```json
{
    "status":"success",
}
```