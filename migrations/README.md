# Migrations

Database migrations for this project live in **`internal/migrations/`**, not here.

- The app uses **golang-migrate** and runs migrations on startup when you run `go run ./cmd/server` (using `DATABASE_URL` from `.env`).
- Rollback: `go run ./cmd/rollback` (rolls back the last migration).

See **PRODUCTION_MIGRATION_GUIDE.md** in the project root for production and staging steps.
