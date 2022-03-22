# scimfe
This repo deploys and configures the PostgreSQL database and Redis instance to support [SCIMPLISTIC](https://github.com/strick-j/scimplistic)

## Prerequisites

* Docker
* docker-compose
* **GNU** Make
* Go 1.12+
    * *Optionals*:
    * [golang-migrate](https://github.com/golang-migrate/migrate) (for manual sql migration)
    * [golangci-lint](https://golangci-lint.run/) as linter

## Tests

### End-to-end

Integration tests (e2e) are described in [e2e](test/e2e/) directory.
Tests cover all cases, including data structure and routes validation.

* Start environment with `docker-compose start`
* Start back-end API with `make run`
* Run tests with `make e2e`

## Usage

### Development
* Start DB and Redis containers with `docker-compose start`located in [deployments](/deployments/)
  * Pre-create containers before start using `docker-compose up -d` (one time operation)
* `make run`

### Production

Use `make` to build the project.
Output binary will be located at `target` directory.

#### Configuration
The service can be configured using environment variables, or a [config file](/configs/)

Use `-c` flag to provide path to a config file.

### Environment Variables
See [config.go](/internal/config/config.go) for more options.

| Name                    | Type   | Defaults                           | Description                                      |
|-------------------------|--------|------------------------------------|--------------------------------------------------|
| `SCIMFE_HTTP_ADDR`      | string | `:8800`                            | Interface to listen by HTTP server               |
| `SCIMFE_DB_ADDRESS`     | string | `postgres://localhost:5432/ledger` | Postgres DB address (URL or DSN)                 |
| `SCIMFE_REDIS_ADDRESS`  | string | `localhost:6379`                   | Redis server address                             |
| `SCIMFE_REDIS_USER`     | string | -                                  | Redis username                                   |
| `SCIMFE_REDIS_PASSWORD` | string | -                                  | Redis password                                   |
| `SCIMFE_REDIS_DB`       | int    | -                                  | Redis database number                            |
| `SCIMFE_MIGRATIONS_DIR` | string | `db/migrations`                    | Path to directory containing migration scripts   |
| `SCIMFE_VERSION_TABLE`  | string | `schema_migrations`                | Name of a table, which contains database version |
| `SCIMFE_SCHEMA_VERSION` | int    | -                                  | Force set schema version (dangerous)             |
| `SCIMFE_NO_MIGRATION`   | bool   | `false`                            | Skip database migration                          |