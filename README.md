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

## Auth

`POST /auth/sign-in`

Params (json):
* email [string]
* password [string]

#### Request
```bash
curl -L -X POST 'localhost:8080/auth/sign-in' -H 'Content-Type: application/json' --data-raw '{
    "email": "test1234test@mmail.com",
    "password": "1234567891"
}'
```

#### Response

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjc1NjgyMjAsImlhdCI6MTYyNzU2NzMyMCwidXNlcl9pZCI6MywidXVpZCI6ImFlZTg3MThkLWRkZjYtNGYwMy05OGM3LTg2ZmE1NGI2MDQyNCJ9.NukptPP9lLLqz29_M0d-lUZeLFk7Tetze3vMhyFrCfQ",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjgxNzIxMjAsImlhdCI6MTYyNzU2NzMyMCwidXNlcl9pZCI6MywidXVpZCI6ImFlZTg3MThkLWRkZjYtNGYwMy05OGM3LTg2ZmE1NGI2MDQyNCJ9.YqimvrtWcx7Vq0ULLFW3H3Lhcov1mZ-b3Kjr9w0x-z4"
}
```