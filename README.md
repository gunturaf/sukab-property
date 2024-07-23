# Sukab Property 

The name Sukab is taken from the fictional character made by [famous Indonesian Poet, Seno G. Ajidarma](https://en.wikipedia.org/wiki/Seno_Gumira_Ajidarma).

## How to Run (development mode)

1. Download and Install Go https://go.dev/dl/ (alternatively, there's also Homebrew for that https://formulae.brew.sh/formula/go)
2. Run PostgreSQL server, a minimum version of 14 is required.
3. Create (or use existing) database, then import table schema for the database in the directory `db_schema`.
4. Set env vars on the terminal, this is an example to set env vars, feel free to modify:
   ```
   export SERVER_LISTEN_ADDR=:8080
   export DB_HOST=localhost 
   export DB_PORT=5432
   export DB_USERNAME=sukab 
   export DB_PASSWORD=<redacted>
   export DB_NAME=sukab_property
   ```
5. Run `go run main.go`

## How to Run (production mode)

### On Linux / MacOS

1. Make sure you have Go toolchain installed.
2. Run `go build -o server-app .`, an executable named `server-app` will be created.
3. Run PostgreSQL server, a minimum version of 14 is required.
4. Create (or use existing) database, then import table schema for the database in the directory `db_schema`.
5. Set env vars on the terminal, this is an example to set env vars, feel free to modify:
   ```
   export SERVER_LISTEN_ADDR=:8080
   export DB_HOST=localhost 
   export DB_PORT=5432
   export DB_USERNAME=sukab 
   export DB_PASSWORD=<redacted>
   export DB_NAME=sukab_property
   ```
6. Run: `./server-app`.

### On Windows

1. Make sure you have Go toolchain installed.
2. Run `go build -o server-app.exe .`, an executable named `server-app.exe` will be created.
3. Run PostgreSQL server, a minimum version of 14 is required.
4. Create (or use existing) database, then import table schema for the database in the directory `db_schema`.
5. Set env vars on the terminal, refer to the one on Linux above.
6. Run the created exe.


## Appendix 1: Environment Variables 

| Key | Description | Default Value |
|-----|-------------|---------------|
| `SERVER_LISTEN_ADDR` | host:port which the HTTP should listen to. | <empty-string> |
| `DB_HOST` | Host for the database, e.g.: "localhost" | <empty-string> |
| `DB_PORT` | Port for the database, e.g.: "5432" | <empty-string> |
| `DB_USERNAME` | Username for the database, e.g.: "postgres" | <empty-string> |
| `DB_PASSWORD` | Password for the database User. | <empty-string> |
| `DB_NAME` | Name for the database, e.g.: "sukab_property". | <empty-string> |
| `DB_ENABLE_SSL` | Enable SSL connection for the database. SSL by default is disabled for local dev. Set empty string to disable. | <empty-string> |

