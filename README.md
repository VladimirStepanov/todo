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

## Refresh JWT token

`POST /auth/refresh`

Params (json):
* refresh_token [string]

#### Request
```bash
curl -L -X POST 'localhost:8080/auth/refresh' -H 'Content-Type: application/json' --data-raw '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjgzMjg0OTEsImlhdCI6MTYyNzcyMzY5MSwidXNlcl9pZCI6MywidXVpZCI6Ijc4YjE2NWE3LWI5MWQtNDM5ZS04MTI4LTBiNmM0ODk3YzNlZSJ9.MRrBSA-T5fS7t249K6-4k70WbUS8--9ZBdA44RybakM"
}'
```

#### Response

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjc1NjgyMjAsImlhdCI6MTYyNzU2NzMyMCwidXNlcl9pZCI6MywidXVpZCI6ImFlZTg3MThkLWRkZjYtNGYwMy05OGM3LTg2ZmE1NGI2MDQyNCJ9.NukptPP9lLLqz29_M0d-lUZeLFk7Tetze3vMhyFrCfQ",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjgxNzIxMjAsImlhdCI6MTYyNzU2NzMyMCwidXNlcl9pZCI6MywidXVpZCI6ImFlZTg3MThkLWRkZjYtNGYwMy05OGM3LTg2ZmE1NGI2MDQyNCJ9.YqimvrtWcx7Vq0ULLFW3H3Lhcov1mZ-b3Kjr9w0x-z4"
}
```

## Logout

`GET /auth/logout`

Params (HTTP header):
* Authorization: Bearer <access_token>

#### Request
```bash
curl -H "Authorization: Bearer ${ACCESS_TOKEN}" localhost:8080/auth/logout
```

#### Response

```json
{"status": "success"}
```

## Create list

`POST /api/lists/`

Params (HTTP header):
* Authorization: Bearer <access_token>

Params (json):
* title [string]
* description [string]


#### Request
```bash
curl -L -X POST 'localhost:8080/api/lists' -H 'Content-Type: application/json' -H "Authorization: Bearer ${ACCESS_TOKEN}" --data-raw '{
    "title": "hello",
    "description": "world"
}'
```

#### Response

```json
{
    "status": "success",
    "list_id": 4
}
```


## Get list by id

`GET /api/lists/:list_id`

Params (HTTP header):
* Authorization: Bearer <access_token>

URL:

* list id (integer)

#### Request
```bash
curl -H "Authorization: Bearer ${ACCESS_TOKEN}" localhost:8080/api/lists/1
```

#### Response

```json
{
    "list_id": 1,
    "title": "title",
    "description": "description"
}
```

## Edit role

`POST /api/lists/:list_id/edit-role`

Params (HTTP header):
* Authorization: Bearer <access_token>

Params (json):
* user_id [int]
* is_admin [bool]

Params (url):

* list id [int]



#### Request
```bash
curl -L -X POST 'localhost:8080/api/lists/:list_id/edit-role' -H 'Content-Type: application/json' -H "Authorization: Bearer ${ACCESS_TOKEN}" --data-raw '{
    "user_id": 10,
    "is_admin": true
}'
```

#### Response

```json
{
    "status": "success"
}
```

## Delete

`POST /api/lists/:list_id/delete`

Params (HTTP header):
* Authorization: Bearer <access_token>

Params (url):

* list id [int]



#### Request
```bash
curl -L -X POST 'localhost:8080/api/lists/:list_id/delete' -H 'Content-Type: application/json' -H "Authorization: Bearer ${ACCESS_TOKEN}"

#### Response

```json
{
    "status": "success"
}
```