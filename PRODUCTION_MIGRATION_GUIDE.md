# Production migration guide: Old schema → New schema

The schema migration is **000011_old_schema_to_new_schema** in **`internal/migrations/`**. It runs automatically with the rest of your migrations (e.g. when the app starts or when you run migrate up).

---

## What to do **now** (dev / data can be dropped)

1. **Backup** the database (optional but recommended even for dev).
2. Run migrations (uses `DATABASE_URL` from `.env`):
   ```bash
   go run ./cmd/server
   ```
   Or run only migrate up (if you have migrate CLI):
   ```bash
   migrate -path internal/migrations -database "$DATABASE_URL" up
   ```
3. If anything fails, fix and re-run (the up script is written to be re-runnable where possible).

Data in **removed columns** (e.g. `books.original_file_path`, `chapters.chapter_number`, `states`, etc.) is **dropped** and not recoverable after the migration.

---

## What to do in **production** (before running the migration)

### 1. Full backup (mandatory)

- Take a **full DB backup** (e.g. `pg_dump`) and store it somewhere safe.
- Optionally take a **snapshot** of the DB volume/instance if your provider supports it.
- Ensure you can **restore** from this backup (test restore in a separate DB if possible).

### 2. Decide what removed data (if any) you must keep

The migration **drops** these columns/tables and their data:

| Removed | Suggestion |
|--------|------------|
| **states** table | If you need state names/codes for reports or compliance, **export** to CSV/backup table before migration (e.g. `CREATE TABLE states_archive AS SELECT * FROM states;` or export to S3/file). |
| **books**: `original_file_path`, `processed_file_path`, `uploaded_by_user_id`, and many other columns | Migration **copies** `uploaded_by_user_id` → `created_by` and `COALESCE(processed_file_path, original_file_path)` → `file_path`. All other removed book columns are **lost**. If you need any of them (e.g. for audits), **export** to an archive table or file before running. |
| **chapters**: `chapter_number`, `page_start`/`page_end`, verification fields, etc. | Not copied. Export to archive if you need them. |
| **users**: `user_type`, `phone_verified_at` | Not copied. Export if needed. |
| **subjects**: `country_id`, `subject_type`, `sequence_order`, `description`, `created_by` | Not copied. Export if needed. |
| **grades**: `display_order`, `numeric_equivalent`, `academic_stage` | Not copied. Export if needed. |

**Action:** For each of the above, decide: “Do we need this data later?” If yes, run a one-time **export** (e.g. `COPY (...) TO STDOUT WITH CSV HEADER` or `INSERT INTO archive_* SELECT ...`) **before** running the migration.

### 3. Maintenance window and rollback plan

- Run the migration in a **maintenance window** (or during low traffic) so you can fix issues without pressure.
- **Rollback plan:** If something goes wrong, restore from the backup taken in step 1. There is **no “reverse migration” script** that recreates the old schema from the new one; rollback = restore backup.
- If you need to **abort** mid-migration: the script runs in a single **transaction** (`BEGIN` … `COMMIT`). If you don’t run `COMMIT` (e.g. you stop the session or run `ROLLBACK`), all changes are rolled back.

### 4. Test on staging first

- Restore a **copy of production** (or a recent backup) to a **staging** DB.
- Run the same migrations on staging (e.g. `go run ./cmd/server` with staging `DATABASE_URL`, or `migrate -path internal/migrations -database "$DATABASE_URL" up`).
- Run your app and smoke tests against the new schema.
- Only then run the migration on production.

### 5. Run the migration on production

- Run migrations (same as staging): start the app so it runs migrate up, or run the migrate CLI with production `DATABASE_URL`:
  ```bash
  # Option A: deploy app (it runs migrations on startup)
  go run ./cmd/server   # or your deploy command; ensure DATABASE_URL is prod DB

  # Option B: migrate CLI
  migrate -path internal/migrations -database "$DATABASE_URL" up
  ```
- If it completes without error, the new schema is in place.
- **Deploy application code** that matches the new schema (e.g. uses `books.created_by`, `user_default`, no `states`, etc.). Old code that still references dropped columns/tables will break.

### 6. After migration

- Smoke-test critical flows (login, subject/book/chapter listing, uploads).
- Keep the backup for at least a few days (or per your retention policy) so you can restore if a problem appears later.

---

## Short checklist (production)

- [ ] Full DB backup taken and restore tested
- [ ] Any “must keep” removed data exported/archived
- [ ] Migration tested on staging
- [ ] Maintenance window communicated
- [ ] Migration run on production
- [ ] App deployed to match new schema
- [ ] Smoke tests passed
- [ ] Backup retained for rollback safety

---

## Summary

| Environment | Data in removed columns | What to do |
|-------------|--------------------------|------------|
| **Dev / now** | OK to delete | Run migration; no extra steps if you don’t need old data. |
| **Production** | May need to keep some | 1) Backup. 2) Export any data you must keep. 3) Test on staging. 4) Run migration in a window. 5) Deploy app. 6) Verify. |

The **same migration** (`internal/migrations/000011_old_schema_to_new_schema.up.sql`) runs in both dev and production (and staging); the difference is the **preparation and rollback plan** in production.
