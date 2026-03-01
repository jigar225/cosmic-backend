-- Migration: Old schema (2026-02-23) → New schema (2026-02-28)
-- Data in removed columns will be dropped (see PRODUCTION_MIGRATION_GUIDE.md for prod).

-- =============================================================================
-- 1. DROP FOREIGN KEYS that reference columns/tables we are removing
-- =============================================================================

-- Drop FKs that reference states (PostgreSQL may name them _fkey or _foreign)
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_context') THEN
    ALTER TABLE user_context DROP CONSTRAINT IF EXISTS user_context_current_state_id_fkey;
    ALTER TABLE user_context DROP CONSTRAINT IF EXISTS user_context_current_state_id_foreign;
    ALTER TABLE user_context DROP CONSTRAINT IF EXISTS user_context_current_subject_id_fkey;
    ALTER TABLE user_context DROP CONSTRAINT IF EXISTS user_context_current_subject_id_foreign;
  ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_default') THEN
    ALTER TABLE user_default DROP CONSTRAINT IF EXISTS user_context_current_state_id_fkey;
    ALTER TABLE user_default DROP CONSTRAINT IF EXISTS user_context_current_state_id_foreign;
    ALTER TABLE user_default DROP CONSTRAINT IF EXISTS user_context_current_subject_id_fkey;
    ALTER TABLE user_default DROP CONSTRAINT IF EXISTS user_context_current_subject_id_foreign;
  END IF;
END $$;

ALTER TABLE generated_content DROP CONSTRAINT IF EXISTS generated_content_state_id_fkey;
ALTER TABLE generated_content DROP CONSTRAINT IF EXISTS generated_content_state_id_foreign;

ALTER TABLE chapters DROP CONSTRAINT IF EXISTS chapters_parent_chapter_id_foreign;
ALTER TABLE chapters DROP CONSTRAINT IF EXISTS chapters_previous_chapter_id_foreign;
ALTER TABLE chapters DROP CONSTRAINT IF EXISTS chapters_next_chapter_id_foreign;
ALTER TABLE chapters DROP CONSTRAINT IF EXISTS chapters_concept_verified_by_foreign;

ALTER TABLE books DROP CONSTRAINT IF EXISTS books_uploaded_by_user_id_foreign;
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_medium_id_foreign;
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_grade_id_foreign;

ALTER TABLE subjects DROP CONSTRAINT IF EXISTS subjects_country_id_foreign;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'states') THEN
    ALTER TABLE states DROP CONSTRAINT IF EXISTS states_country_id_foreign;
    ALTER TABLE states DROP CONSTRAINT IF EXISTS states_default_board_id_foreign;
  END IF;
END $$;

DROP TABLE IF EXISTS states;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_context') THEN
    ALTER TABLE user_context RENAME TO user_default;
  END IF;
END $$;
ALTER TABLE user_default DROP COLUMN IF EXISTS current_subject_id;

ALTER TABLE countries DROP COLUMN IF EXISTS updated_at;
ALTER TABLE countries DROP COLUMN IF EXISTS has_states;

ALTER TABLE boards DROP COLUMN IF EXISTS updated_at;

DROP INDEX IF EXISTS grades_numeric_equivalent_index;
ALTER TABLE grades DROP COLUMN IF EXISTS display_order;
ALTER TABLE grades DROP COLUMN IF EXISTS numeric_equivalent;
ALTER TABLE grades DROP COLUMN IF EXISTS academic_stage;

ALTER TABLE mediums DROP COLUMN IF EXISTS script_type;

DROP INDEX IF EXISTS subjects_country_id_board_id_medium_id_grade_id_index;
ALTER TABLE subjects DROP COLUMN IF EXISTS country_id;
ALTER TABLE subjects DROP COLUMN IF EXISTS subject_type;
ALTER TABLE subjects DROP COLUMN IF EXISTS sequence_order;
ALTER TABLE subjects DROP COLUMN IF EXISTS description;
ALTER TABLE subjects DROP COLUMN IF EXISTS created_by;

ALTER TABLE users DROP COLUMN IF EXISTS user_type;
ALTER TABLE users DROP COLUMN IF EXISTS phone_verified_at;
ALTER TABLE users ADD COLUMN IF NOT EXISTS preferable_subject VARCHAR(255)[] NOT NULL DEFAULT '{}';
ALTER TABLE users ADD COLUMN IF NOT EXISTS plateform_version VARCHAR(255) NOT NULL DEFAULT '0.0.0';

ALTER TABLE books ADD COLUMN IF NOT EXISTS created_by BIGINT NULL;
ALTER TABLE books ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE books ADD COLUMN IF NOT EXISTS file_path VARCHAR(1000) NULL;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'books' AND column_name = 'uploaded_by_user_id') THEN
    UPDATE books SET created_by = uploaded_by_user_id WHERE uploaded_by_user_id IS NOT NULL;
  END IF;
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'books' AND column_name = 'original_file_path') THEN
    UPDATE books SET file_path = COALESCE(processed_file_path, original_file_path) WHERE file_path IS NULL;
  ELSIF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'books' AND column_name = 'processed_file_path') THEN
    UPDATE books SET file_path = processed_file_path WHERE file_path IS NULL AND processed_file_path IS NOT NULL;
  END IF;
END $$;

UPDATE books SET file_path = '' WHERE file_path IS NULL;
ALTER TABLE books ALTER COLUMN file_path SET NOT NULL;

DROP INDEX IF EXISTS books_subject_id_medium_id_grade_id_index;
DROP INDEX IF EXISTS books_book_type_index;
DROP INDEX IF EXISTS books_status_index;

ALTER TABLE books DROP COLUMN IF EXISTS book_type;
ALTER TABLE books DROP COLUMN IF EXISTS medium_id;
ALTER TABLE books DROP COLUMN IF EXISTS grade_id;
ALTER TABLE books DROP COLUMN IF EXISTS author;
ALTER TABLE books DROP COLUMN IF EXISTS edition;
ALTER TABLE books DROP COLUMN IF EXISTS isbn;
ALTER TABLE books DROP COLUMN IF EXISTS book_code;
ALTER TABLE books DROP COLUMN IF EXISTS uploaded_by_user_id;
ALTER TABLE books DROP COLUMN IF EXISTS curriculum_version;
ALTER TABLE books DROP COLUMN IF EXISTS effective_start_date;
ALTER TABLE books DROP COLUMN IF EXISTS effective_end_date;
ALTER TABLE books DROP COLUMN IF EXISTS status;
ALTER TABLE books DROP COLUMN IF EXISTS original_file_path;
ALTER TABLE books DROP COLUMN IF EXISTS processed_file_path;
ALTER TABLE books DROP COLUMN IF EXISTS cover_image_url;
ALTER TABLE books DROP COLUMN IF EXISTS has_toc;
ALTER TABLE books DROP COLUMN IF EXISTS toc_extraction_method;
ALTER TABLE books DROP COLUMN IF EXISTS processing_status;
ALTER TABLE books DROP COLUMN IF EXISTS processing_notes;
ALTER TABLE books DROP COLUMN IF EXISTS processing_started_at;
ALTER TABLE books DROP COLUMN IF EXISTS processing_completed_at;
ALTER TABLE books DROP COLUMN IF EXISTS description;
ALTER TABLE books DROP COLUMN IF EXISTS tags;
ALTER TABLE books DROP COLUMN IF EXISTS view_count;
ALTER TABLE books DROP COLUMN IF EXISTS download_count;
ALTER TABLE books DROP COLUMN IF EXISTS updated_at;

CREATE INDEX IF NOT EXISTS books_is_active_index ON books(is_active);
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'books_created_by_foreign') THEN
    ALTER TABLE books ADD CONSTRAINT books_created_by_foreign FOREIGN KEY (created_by) REFERENCES users(id);
  END IF;
END $$;

DROP INDEX IF EXISTS chapters_book_id_chapter_number_index;
DROP INDEX IF EXISTS chapters_status_index;

ALTER TABLE chapters DROP COLUMN IF EXISTS chapter_number;
ALTER TABLE chapters DROP COLUMN IF EXISTS chapter_code;
ALTER TABLE chapters DROP COLUMN IF EXISTS page_start;
ALTER TABLE chapters DROP COLUMN IF EXISTS page_end;
ALTER TABLE chapters DROP COLUMN IF EXISTS section_type;
ALTER TABLE chapters DROP COLUMN IF EXISTS word_count;
ALTER TABLE chapters DROP COLUMN IF EXISTS learning_objectives;
ALTER TABLE chapters DROP COLUMN IF EXISTS key_concepts;
ALTER TABLE chapters DROP COLUMN IF EXISTS estimated_teaching_hours;
ALTER TABLE chapters DROP COLUMN IF EXISTS difficulty_level;
ALTER TABLE chapters DROP COLUMN IF EXISTS prerequisites;
ALTER TABLE chapters DROP COLUMN IF EXISTS concept_extraction_confidence;
ALTER TABLE chapters DROP COLUMN IF EXISTS concept_extraction_method;
ALTER TABLE chapters DROP COLUMN IF EXISTS concept_verified;
ALTER TABLE chapters DROP COLUMN IF EXISTS concept_verified_by;
ALTER TABLE chapters DROP COLUMN IF EXISTS concept_verified_at;
ALTER TABLE chapters DROP COLUMN IF EXISTS embedding_generated_at;
ALTER TABLE chapters DROP COLUMN IF EXISTS embedding_model;
ALTER TABLE chapters DROP COLUMN IF EXISTS embedding_dimensions;
ALTER TABLE chapters DROP COLUMN IF EXISTS status;
ALTER TABLE chapters DROP COLUMN IF EXISTS version;
ALTER TABLE chapters DROP COLUMN IF EXISTS parent_chapter_id;
ALTER TABLE chapters DROP COLUMN IF EXISTS previous_chapter_id;
ALTER TABLE chapters DROP COLUMN IF EXISTS next_chapter_id;
ALTER TABLE chapters DROP COLUMN IF EXISTS view_count;
ALTER TABLE chapters DROP COLUMN IF EXISTS updated_at;
