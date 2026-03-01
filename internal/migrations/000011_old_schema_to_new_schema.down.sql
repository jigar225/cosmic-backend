-- Rollback for 000011 (old schema → new schema) is not supported.
-- This migration drops columns and the states table; reversing would require
-- full schema and data restoration. Use a DB backup to roll back if needed.
-- See PRODUCTION_MIGRATION_GUIDE.md.
SELECT 1;
