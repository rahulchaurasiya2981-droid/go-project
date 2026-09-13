1. Setup the app enviroment (development,production)
2. Load env configuration file (.env) in app
3. Setup the custome logger for logs (json)
4.

# Flow of `go run .\cmd\api`

- `go run .\cmd\api`
- Run will compile the main.go inside api folder into maschine code (api.exe)
- and this api.exe is loade in memory by os to execute

```go
                    You
                     │
                     │  go run .\cmd\api
                     ▼
              ┌───────────────┐
              │    Go CLI     │
              │   `go run`    │
              └───────┬───────┘
                      │
                      ▼
            Find package in
              .\cmd\api
                      │
                      ▼
              ┌───────────────┐
              │ Go Compiler   │
              │    (compile)  │
              └───────┬───────┘
                      │
                      │ Go source code
                      │ → native machine code
                      ▼
              ┌────────────────┐
              │ Temporary      │
              │ executable     │
              │    api.exe     │
              └───────┬────────┘
                      │
                      │ OS starts process
                      ▼
              ┌────────────────┐
              │    Windows     │
              │       OS       │
              └───────┬────────┘
                      │
                      │ Loads executable
                      │ into process memory
                      ▼
          ┌────────────────────────┐
          │       Process          │
          │                        │
          │  Code  → Memory        │
          │  Data  → Memory        │
          │  Heap  → Memory        │
          │  Stack → Memory        │
          └───────────┬────────────┘
                      │
                      ▼
                main.main()
                      │
                      ▼
               Your API starts
                      │
                      ▼
              ┌───────────────┐
              │ HTTP Server   │
              │   :8080       │
              └───────────────┘
```

# go clean -cache

- `go clean -cache`
- Delete Go's compiled build cache.

```go
    go clean -cache
        │
        ▼
    Go finds its build cache
        │
        ▼
    Deletes cached build artifacts
        │
        ▼
    Cache is empty
```

- Go has several types of cached/downloaded data :

```go
Go
│
├── Build Cache
│     └── go clean -cache
│
├── Module Download Cache
│     └── go clean -modcache
│
└── Test Cache
      └── go clean -testcache
```

# API Backend Flow

```go
      Request
      ↓
      Route
      ↓
      Handler
      ↓
      Service
      ↓
      Repository
      ↓
      database/sql
      ↓
      PostgreSQL
```

| Layer      | Responsibility        |
| ---------- | --------------------- |
| Route      | URL + HTTP method     |
| Handler    | HTTP request/response |
| Service    | Business logic        |
| Repository | Database queries      |
| Model      | Data structures       |
| Database   | PostgreSQL connection |

# Makefile in Go

- A Makefile is an automation script tool (part of the Unix make build system) used in Go projects to simplify long, repetitive terminal commands into short, memorable shortcuts.
- In Go, commands like running tests, building binaries, applying database migrations, or running linters can get long and difficult to remember. A Makefile acts as a single control panel for your project.

# Database migrations

---

- It is used the creat or intialize the tables in Database (if table not exist it will create)
- go long support database migration with using below librarby

```go
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
```

- Steps to create migration

```go
Step 0 : Start : 'main()'
Step 1 : `Check` is migration command is provided (`len(os.Args) < 2`)
Step 2 : `Load` configuratio file (that load `.env`)
Step 3 : `Initialize` migration instance (using database connection url)
Step 4 : `Extract` migration command `command := os.Args[1]`
Step 4 : `Handle` migration command `(up/down)` by 'switch command{}'.
Step 5 : End : 'main()'
```

- After running migration up it will create 2 tables in Database (crud_db)

```go
1. 'users' (store business data)
2. 'schema_migrations' (auto create and store migration state)
```

- `schema_migrations` is infrastructure used by `golang-migrate` library for tracking purpose.
- Think of schema_migrations as the `memory/bookkeeping system of golang-migrate`

```go
                 'crud_db'
                    │
          ┌─────────┴─────────┐
          │                   │
          ▼                   ▼
       'users'          'schema_migrations'
          │                   │
          │                   └── `golang-migrate's migration state`
          │
          └── `your application data`
```

## How to check migration detils in DB

- `version` = 1 → migration 001_create_users has been applied.
- `dirty` = false → migration completed successfully.

```sql
crud_db=# SELECT * FROM schema_migrations;
 version | dirty
---------+-------
       2 | f
(1 row)
```

```go
000001_create_users_table.up.sql
              │
              └── Apply migration


000001_create_users_table.down.sql
              │
              └── Rollback migration
```

# Project Flow

## `Step 1` → PostgreSQL connection

- PostgreSQL is relation database
- You can use this DB by `local installing`,using `docker (*)` and `online service`.
- go lang support this drive for DB : `https://go.dev/wiki/SQLDrivers`
- We will use `http://github.com/jackc/pgx` for postgreSql Database. (modern & updated driver & actively maintained.)
- We can also used directly postgresql driver code and connected with db but what if in future we change the database than we need install their respective driver and again write the that db driver code from scratch again so avoiding this we use interface (`database/sql`) for different driver
- `database/sql` is a common API/abstraction

### Architecture of Backend Database

- database/sql provides a common Go API for working with databases. Database-specific drivers, such as pgx for PostgreSQL, implement the functionality needed by database/sql to communicate with their respective databases.

```go
             Go Application
                   │
                   ▼
             database/sql
          Common database API
                   │
                   ▼
             PostgreSQL Driver
                 pgx
                   │
                   ▼
              PostgreSQL
```

```go
             Go Application
                   │
                   ▼
             database/sql
                   │
          ┌────────┴────────┐
          ▼                 ▼
        pgx driver       MySQL driver
          │                 │
          ▼                 ▼
     PostgreSQL           MySQL
```

- pgx itself has its own native API, but stdlib provides the database/sql compatibility layer/driver for pgx.
- database/sql gives your Go application a standard database interface/API, while the driver adapts that standard API to the specific database. This reduces coupling to a particular database driver and makes future database changes easier.

### `_ "github.com/jackc/pgx/v5/stdlib"`

- `Load` pgx's database/sql driver and run its `initialization`, but don't give me a package name to use in this file.
- Database/sql itself doesn't know how to communicate with PostgreSQL.

```go
import "database/sql"
```

- The "pgx" driver name works because pgx/stdlib registered it during package initialization.

```go
db, err := sql.Open("pgx", dsn)
```

## `Step 2` → Test DB connection

- run the project

```go
{
      "time":"2026-09-12T20:07:11.0016205+05:30",
      "level":"INFO",
      "msg":"Database connection established successfully",
      "action":"DB_CONNECTED",
      "host":"localhost",
      "database":"crud_db"
}
```

## `Step 3` → Create model/table

`Step 4` → Repository layer
`Step 5` → Service layer
`Step 6` → CRUD handlers
`Step 7` → Routes
`Step 8` → Postman testing

# Notes :

- godotenv.Load() -> load bydefault .env file
