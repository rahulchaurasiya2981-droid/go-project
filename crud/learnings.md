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

- Database migrations is management of incremental,reversible changes to relational database schemas.
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

## Database Migration Rollback Guard

---

- m.Down() : will rollback enitre migration (means drop all tables in database)
- is very dangerous to drop all table from production database
- so we use m.Step(-1)`

# Prevent Zombie Queries using Context

- Create the issues
  Step 1 : add pg_sleep(20) in query
  Step 2 : Hit from server and cencel the request
  Step 3 : run below command and see still client request is executing in database but client already cancel the request
- this will create server crash issues if their is lot of clients

## command to check list of active query in Database

```sql
SELECT
    pid,
    state,
    now() - query_start AS duration,
    wait_event,
    query
FROM pg_stat_activity
WHERE query ILIKE '%pg_sleep%';
```

- Go give feature to use `context` than this signal is passed to query (client is already cancel the request so stop this query which are running in background)

## User context

```go

```

```go
rows, err := db.Query(query)
```

# Folder Architecture in GO

1. Package-by-Feature / Domain-Driven (Recommended)

```go
project-root/
├── cmd/
│   └── api/
│       └── main.go           // Entry point & dependency injection
├── internal/
│   ├── user/                // Everything user-related lives here
│   │   ├── entity.go        // DB Struct / Domain Entity
│   │   ├── dto.go           // Request / Response Structs
│   │   ├── handler.go       // HTTP Handler / Controller
│   │   ├── service.go       // Business Logic
│   │   ├── repository.go    // DB Queries
│   │   └── router.go        // HTTP Endpoints
│   ├── order/               // Everything order-related lives here
│   │   ├── entity.go
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   └── platform/            // Shared cross-cutting concerns
│       ├── database/        // Connection pools (PostgreSQL, Redis)
│       └── logger/          // Logging setup (slog/zap)
├── pkg/                     // Generic utility packages (JWT, Crypto)
├── go.mod
└── go.sum
```

2. Layered Architecture (Package-by-Layer)

```go
project-root/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/            // All HTTP handlers
│   │   ├── user.go
│   │   └── order.go
│   ├── services/            // All business logic
│   │   ├── user.go
│   │   └── order.go
│   ├── repositories/        // All database queries
│   │   ├── user.go
│   │   └── order.go
│   └── models/              // All shared structs & DTOs
│       ├── user.go
│       └── order.go
├── go.mod
└── go.sum
```

3. Flat / Single-Package Structure
4. Clean / Hexagonal Architecture (Ports & Adapters)

| Architecture Style     | Refactoring Effort | Circular Import Risk | Recommended Team Size     |
| ---------------------- | ------------------ | -------------------- | ------------------------- |
| **Flat Structure**     | Hard               | None                 | 1 Developer               |
| **Layered Structure**  | Medium             | High                 | Small Teams (1–3)         |
| **Package-by-Feature** | Easy               | Very Low             | Small to Large Teams (3+) |
| **Clean / Hexagonal**  | Medium             | None                 | Enterprise Teams          |

- In this project we using `Package-by-Feature`

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
