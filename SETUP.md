# Go + PostgreSQL + Fiber – Database and server setup (step by step)

Follow these steps in order. This is a **generic template** for any Go app that uses:

- PostgreSQL
- A `migrations/` folder with SQL files
- A Fiber HTTP server (`github.com/gofiber/fiber/v2`) started from `cmd/server`
- **Split transport layer**: `internal/transport/app.go` (Fiber app creation), `internal/transport/routes.go` (all routes), `internal/transport/handlers/` (handler files per concern)—**not** one big `server.go`

Assumption: you already have PostgreSQL running on your laptop.

---

## Step 0: Environment file (`.env`)

Copy the example env file (if present) to `.env` and fill in the values you need:

- **Database**: `DATABASE_URL` if you are not using the default local Postgres (`postgres/postgres` on `localhost:5432` with your app database).
- **S3 (optional)**: `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `S3_BUCKET` if you want image upload / board image URLs to work.

The server automatically loads `.env` on startup so `os.Getenv` can read these values.

---

## Step 1: Create the database (you do this manually)

PostgreSQL must be running. **Create the app database yourself** using whatever you prefer:

- pgAdmin / TablePlus / DBeaver UI  
- `psql` in your terminal  
- Cloud console (RDS, Cloud SQL, etc.)

Name it something like `app_db_name`.  
Once the DB exists, the Go app will use migrations to create/update tables.

---

## Step 2: Set connection (optional)

The app connects to the database you created in Step 1.  
By default, it expects something like:

- **Host:** localhost  
- **Port:** 5432  
- **User:** postgres  
- **Password:** postgres  
- **Database:** `app_db_name` (the database you created)  

If that matches your setup, you don’t need to do anything.

If that does **not** match your setup, set `DATABASE_URL` before running the server (replace placeholders with your own values):

```bash
export DATABASE_URL="postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/app_db_name?sslmode=disable"
```

Replace `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, and `app_db_name` with your PostgreSQL values.

---

## Step 3: Run the Fiber server (from project root)

Open a terminal, go to the project root (where `go.mod` and `migrations/` are), then run:

```bash
cd /path/to/your/project-root
go run ./cmd/server
```

You should see logs similar to:

- `Migrations applied.`
- `Server started. DB connected. Fiber HTTP server on :8080`

Migrations run automatically on startup and create/update your app tables (for example, a `boards` table).

---

## Step 4: Test the API (Fiber HTTP server)

In another terminal (leave the server running):

**Health check (Fiber route)**

```bash
curl http://localhost:8080/health
```

Expected: `{"status":"ok"}`

**Example resource check (adjust to your app)**

If your app exposes a simple resource like `boards`, you might have:

```bash
curl http://localhost:8080/boards
```

Expected at first: `[]`

```bash
curl -X POST http://localhost:8080/boards \
  -H "Content-Type: application/json" \
  -d '{"name":"Example"}'
```

Expected: JSON with fields like `id`, `name`, `created_at`.

---

## What each part does (code structure template)

| Part | Role |
|------|------|
| **domain** | Core types / entities (e.g. `Board`, `User`). No DB, HTTP, or S3 calls here. |
| **repository** | Connects to Postgres, runs migrations, runs SQL (create/get/list/update entities). |
| **transport** | Fiber HTTP layer. **Must be split** (see layout below)—do **not** put everything in one `server.go`. |
| **storage** | S3 (or other object storage) integration, wired from env and passed into handlers / services. |
| **service** | Higher‑level business logic that can sit between transport and repository (e.g. validation, workflows). |
| **ingest** | For ingestion pipelines, batch jobs, background processing (optional). |
| **cmd/server** | Main entrypoint. Connects to DB, runs migrations, configures storage, starts Fiber HTTP server, graceful shutdown. |
| **cmd/rollback** | CLI entrypoint to run down/rollback migrations if needed. |

Data flow template: **HTTP (Fiber) → transport → (service) → repository → PostgreSQL**, with optional S3/object storage via `storage` when uploading or serving files.

### Transport layout (required – avoid one big server.go)

Keep `internal/transport` **split** so the HTTP layer stays maintainable:

| File / folder | Purpose |
|---------------|---------|
| **`internal/transport/app.go`** | Creates the Fiber app (`NewApp`), wires handlers and routes. No route definitions here. |
| **`internal/transport/routes.go`** | **Single place** for all route registration: `RegisterRoutes(app, h)`. One place to see the full API surface. |
| **`internal/transport/handlers/`** | Handler implementations. One file per concern (e.g. `health.go`, `boards.go`, `upload.go`) plus `handlers.go` for the `Handlers` struct that holds `DB`, `Uploader`, etc. |

Do **not** create one `server.go` that mixes app creation, route registration, and handler logic—that leads to a single huge file. Use the structure above.

---

## If something fails

- **“connection refused”**  
  PostgreSQL is not running or not on `localhost:5432`. Start the Postgres server and check host/port.

- **“database app_db_name does not exist”**  
  Make sure you actually created the database yourself in Step 1, and that `DATABASE_URL` (or default config) points to the correct name.

- **“password authentication failed”**  
  Set `DATABASE_URL` (Step 2) with the correct user and password.

- **“migrate: ..."**  
  Ensure you run `go run ./cmd/server` from the project root so the `migrations` folder is found.
  
Once these steps work, you have a Go + Postgres + Fiber app connected to your database with basic CRUD behaviour; you can then add more entities, routes, and business logic as needed.
