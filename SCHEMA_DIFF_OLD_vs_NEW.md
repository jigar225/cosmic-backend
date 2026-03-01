# Schema comparison: Old (2026-02-23) vs New (2026-02-28)

Per-table differences. **Old** = `drawSQL-pgsql-export-2026-02-23 (2).sql`, **New** = `drawSQL-pgsql-export-2026-02-28.sql`.

---

## 1. `countries`

| Change | Detail |
|--------|--------|
| **Removed column** | `updated_at` (TIMESTAMP WITH TIME ZONE) |
| **Removed column** | `has_states` (BOOLEAN NOT NULL DEFAULT FALSE) |

All other columns unchanged.

---

## 2. `grade_methods`

No column or index changes. Same in both schemas.

---

## 3. `boards`

| Change | Detail |
|--------|--------|
| **Removed column** | `updated_at` (TIMESTAMP WITH TIME ZONE) |

Indexes and PK unchanged.

---

## 4. `states` (table)

| Change | Detail |
|--------|--------|
| **Table removed** | Entire table `states` does not exist in new schema. All FKs and references to `states` are gone. |

Old had: `id`, `country_id`, `code`, `title`, `default_board_id`, `is_visible`, `created_at`, `updated_at`, indexes, FKs from `user_context` and `generated_content`.

---

## 5. `grades`

| Change | Detail |
|--------|--------|
| **Removed column** | `display_order` (INTEGER NOT NULL) |
| **Removed column** | `numeric_equivalent` (INTEGER NULL) |
| **Removed column** | `academic_stage` (VARCHAR(50) NULL) |
| **Removed index** | `grades_numeric_equivalent_index` |

All other columns and PK/index unchanged.

---

## 6. `mediums`

| Change | Detail |
|--------|--------|
| **Removed column** | `script_type` (VARCHAR(50) NULL) |

All other columns, indexes, and FKs unchanged.

---

## 7. `subjects`

| Change | Detail |
|--------|--------|
| **Removed column** | `country_id` (BIGINT NOT NULL) |
| **Removed column** | `subject_type` (VARCHAR(50) NOT NULL DEFAULT 'core') |
| **Removed column** | `sequence_order` (INTEGER NULL) |
| **Removed column** | `description` (TEXT NULL) |
| **Removed column** | `created_by` (BIGINT NULL) |
| **Removed index** | `subjects_country_id_board_id_medium_id_grade_id_index` |

FK: old had `subjects_country_id_foreign` → **removed** in new (no `country_id`). All other FKs unchanged.

---

## 8. `users`

| Change | Detail |
|--------|--------|
| **Removed column** | `user_type` (VARCHAR(50) NOT NULL DEFAULT 'individual') |
| **Removed column** | `phone_verified_at` (TIMESTAMP WITH TIME ZONE NULL) |
| **Added column** | `preferable_subject` (VARCHAR(255)[] NOT NULL) |
| **Added column** | `plateform_version` (VARCHAR(255) NOT NULL) |

All other columns and PK unchanged.

---

## 9. `user_context` → `user_default` (rename + structure)

| Change | Detail |
|--------|--------|
| **Table renamed** | `user_context` → `user_default` |
| **Removed column** | `current_subject_id` (BIGINT NULL) |
| **Removed FK** | `user_context_current_subject_id_foreign` |
| **Removed FK** | `user_context_current_state_id_foreign` (states table gone) |

New table has no reference to `states` or `current_subject_id`. All other columns (`user_id`, `current_country_id`, `current_state_id`, `current_board_id`, `current_medium_id`, `current_grade_id`, `updated_at`) and their FKs remain (except state/subject).

---

## 10. `books`

| Change | Detail |
|--------|--------|
| **Removed column** | `book_type` (VARCHAR(50) NOT NULL) |
| **Removed column** | `medium_id` (BIGINT NOT NULL) |
| **Removed column** | `grade_id` (BIGINT NOT NULL) |
| **Removed column** | `author` (VARCHAR(255) NULL) |
| **Removed column** | `edition` (VARCHAR(50) NULL) |
| **Removed column** | `isbn` (VARCHAR(20) NULL) |
| **Removed column** | `book_code` (VARCHAR(100) NULL) |
| **Removed column** | `uploaded_by_user_id` (BIGINT NULL) |
| **Removed column** | `curriculum_version` (VARCHAR(50) NULL) |
| **Removed column** | `effective_start_date` (DATE NULL) |
| **Removed column** | `effective_end_date` (DATE NULL) |
| **Removed column** | `status` (VARCHAR(50) NOT NULL DEFAULT 'draft') |
| **Removed column** | `original_file_path` (VARCHAR(1000) NOT NULL) |
| **Removed column** | `processed_file_path` (VARCHAR(1000) NULL) |
| **Removed column** | `cover_image_url` (VARCHAR(500) NULL) |
| **Removed column** | `has_toc` (BOOLEAN NULL) |
| **Removed column** | `toc_extraction_method` (VARCHAR(50) NULL) |
| **Removed column** | `processing_status` (VARCHAR(50) NULL) |
| **Removed column** | `processing_notes` (TEXT NULL) |
| **Removed column** | `processing_started_at` (TIMESTAMP WITH TIME ZONE NULL) |
| **Removed column** | `processing_completed_at` (TIMESTAMP WITH TIME ZONE NULL) |
| **Removed column** | `description` (TEXT NULL) |
| **Removed column** | `tags` (TEXT[] NULL) |
| **Removed column** | `view_count` (INTEGER NULL) |
| **Removed column** | `download_count` (INTEGER NULL) |
| **Removed column** | `updated_at` (TIMESTAMP WITH TIME ZONE) |
| **Added column** | `created_by` (BIGINT NULL) — replaces uploader concept |
| **Added column** | `is_active` (BOOLEAN NOT NULL DEFAULT 'draft') — note: default 'draft' is odd for BOOLEAN; may be intended as status |
| **Added column** | `file_path` (VARCHAR(1000) NOT NULL) — single path, replaces original/processed |
| **Removed index** | `books_subject_id_medium_id_grade_id_index` |
| **Removed index** | `books_book_type_index` |
| **Removed index** | `books_status_index` |
| **Added index** | `books_is_active_index` |
| **FK change** | `books_uploaded_by_user_id_foreign` → `books_created_by_foreign` (created_by → users.id) |
| **Removed FKs** | `books_medium_id_foreign`, `books_grade_id_foreign` |

---

## 11. `chapters`

| Change | Detail |
|--------|--------|
| **Removed column** | `chapter_number` (INTEGER NOT NULL) |
| **Removed column** | `chapter_code` (VARCHAR(100) NULL) |
| **Removed column** | `page_start` (INTEGER NULL) |
| **Removed column** | `page_end` (INTEGER NULL) |
| **Removed column** | `section_type` (VARCHAR(50) NULL) |
| **Removed column** | `word_count` (INTEGER NULL) |
| **Removed column** | `learning_objectives` (TEXT[] NULL) |
| **Removed column** | `key_concepts` (TEXT[] NULL) |
| **Removed column** | `estimated_teaching_hours` (DECIMAL(4,2) NULL) |
| **Removed column** | `difficulty_level` (VARCHAR(50) NULL) |
| **Removed column** | `prerequisites` (TEXT[] NULL) |
| **Removed column** | `concept_extraction_confidence` (DECIMAL(3,2) NULL) |
| **Removed column** | `concept_extraction_method` (VARCHAR(50) NULL) |
| **Removed column** | `concept_verified` (BOOLEAN NULL) |
| **Removed column** | `concept_verified_by` (BIGINT NULL) |
| **Removed column** | `concept_verified_at` (TIMESTAMP WITH TIME ZONE NULL) |
| **Removed column** | `embedding_generated_at` (TIMESTAMP WITH TIME ZONE NULL) |
| **Removed column** | `embedding_model` (VARCHAR(50) NULL) |
| **Removed column** | `embedding_dimensions` (INTEGER NULL) |
| **Removed column** | `status` (VARCHAR(50) NOT NULL DEFAULT 'draft') |
| **Removed column** | `version` (INTEGER NOT NULL DEFAULT 1) |
| **Removed column** | `parent_chapter_id` (BIGINT NULL) |
| **Removed column** | `previous_chapter_id` (BIGINT NULL) |
| **Removed column** | `next_chapter_id` (BIGINT NULL) |
| **Removed column** | `view_count` (INTEGER NULL) |
| **Removed column** | `updated_at` (TIMESTAMP WITH TIME ZONE) |
| **Removed index** | `chapters_book_id_chapter_number_index` |
| **Removed index** | `chapters_status_index` |
| **Removed FKs** | `chapters_parent_chapter_id_foreign`, `chapters_previous_chapter_id_foreign`, `chapters_next_chapter_id_foreign`, `chapters_concept_verified_by_foreign` |

Chapters in new schema: `id`, `book_id`, `chapter_title`, `content_summary`, `concept_tags`, `embedding_id`, `is_visible`, `created_at` only.

---

## 12. `generated_content`

No column or index differences between old and new. All columns and indexes are the same.  
FKs that referenced removed tables (`states`) still exist in new schema for `state_id` — ensure `states` table exists if you use that FK, or drop the column/FK in a future migration.

---

## Summary: tables

| Table            | Old schema     | New schema     | Notes                    |
|-----------------|----------------|----------------|--------------------------|
| countries       | ✓              | ✓              | Columns simplified       |
| grade_methods   | ✓              | ✓              | Same                     |
| boards          | ✓              | ✓              | No updated_at            |
| states          | ✓              | **Removed**    | Table dropped            |
| grades          | ✓              | ✓              | Fewer columns/index      |
| mediums         | ✓              | ✓              | script_type removed      |
| subjects        | ✓              | ✓              | No country_id, etc.      |
| users           | ✓              | ✓              | New cols, some removed   |
| user_context    | ✓              | —              | Renamed → user_default   |
| user_default    | —              | ✓              | No current_subject_id    |
| books           | ✓              | ✓              | Major simplification     |
| chapters        | ✓              | ✓              | Major simplification     |
| generated_content | ✓            | ✓              | Unchanged                |

---

## Naming / FK consistency notes

1. **books.created_by** (new) vs **books.uploaded_by_user_id** (old): same meaning; new schema uses `created_by` and `books_created_by_foreign`.
2. **books.is_active** in new schema has default `'draft'` in the export — that’s a string; typically this would be BOOLEAN `TRUE`/`FALSE`. Worth fixing in DB (e.g. `DEFAULT TRUE` or a separate `status` column).
3. **user_default** still has **current_state_id** in new schema; table **states** is removed, so that FK will break unless you drop `current_state_id` or re-add `states`.
4. **Orphan state columns:** In the new schema, `user_default.current_state_id` and `generated_content.state_id` still exist, but there is no `states` table and no FKs to it. So those columns are just nullable BIGINTs with no referential integrity. Consider dropping them or re-adding a minimal `states` table if you need state semantics.

Use this diff to write migrations (e.g. add/remove columns, rename table, drop `states`, fix FKs and defaults) when moving from old to new schema.
