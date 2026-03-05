-- Restore grade fields dropped in migration 011, and restore has_states in countries.
-- These are used by the admin panel frontend for curriculum management.

-- Grades: add back display_order, numeric_equivalent, academic_stage
ALTER TABLE grades ADD COLUMN IF NOT EXISTS display_order    INTEGER     NOT NULL DEFAULT 0;
ALTER TABLE grades ADD COLUMN IF NOT EXISTS numeric_equivalent INTEGER   NULL;
ALTER TABLE grades ADD COLUMN IF NOT EXISTS academic_stage   VARCHAR(50) NULL;

-- Countries: add back has_states
ALTER TABLE countries ADD COLUMN IF NOT EXISTS has_states BOOLEAN NOT NULL DEFAULT FALSE;
