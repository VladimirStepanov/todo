# Description

Simple REST API TODO application written using Clean Architecture principles


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
make test
```

# Endpoints

## Registration

`POST /auth/sign-up`

Params (json):
* email [string]
* password [string]

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