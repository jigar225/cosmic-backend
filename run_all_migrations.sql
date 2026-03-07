-- =============================================================================
-- COMBINED MIGRATION — Final schema (equivalent to migrations 000001–000014)
-- Safe to run on a FRESH (empty) database.
-- Safe to re-run on a database where some tables already exist.
-- NOTE: ADD CONSTRAINT IF NOT EXISTS is not valid PostgreSQL syntax.
--       Constraints are wrapped in DO $$ blocks instead.
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ─────────────────────────────────────────────────────────────────────────────
-- countries
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS countries (
    id             BIGSERIAL PRIMARY KEY,
    country_code   VARCHAR(3)   NOT NULL,
    title          VARCHAR(100) NOT NULL,
    icon_path      VARCHAR(255) NULL,
    phone_code     VARCHAR(10)  NULL,
    signup_methods TEXT[]       NOT NULL DEFAULT ARRAY['email']::TEXT[],
    have_board     BOOLEAN      NOT NULL DEFAULT TRUE,
    is_visible     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'countries_country_code_unique') THEN
    ALTER TABLE countries ADD CONSTRAINT countries_country_code_unique UNIQUE (country_code);
  END IF;
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- grade_methods
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS grade_methods (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    description TEXT         NULL,
    is_visible  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'grade_methods_title_unique') THEN
    ALTER TABLE grade_methods ADD CONSTRAINT grade_methods_title_unique UNIQUE (title);
  END IF;
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- languages
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS languages (
    id         BIGSERIAL PRIMARY KEY,
    code       VARCHAR(10)  NOT NULL,
    name       VARCHAR(100) NOT NULL,
    is_visible BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'languages_code_unique') THEN
    ALTER TABLE languages ADD CONSTRAINT languages_code_unique UNIQUE (code);
  END IF;
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- boards
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS boards (
    id              BIGSERIAL PRIMARY KEY,
    country_id      BIGINT       NOT NULL REFERENCES countries(id),
    title           VARCHAR(100) NOT NULL,
    grade_method_id BIGINT       NULL REFERENCES grade_methods(id),
    is_visible      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS boards_country_id_index ON boards(country_id);
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'boards_country_id_title_unique') THEN
    ALTER TABLE boards ADD CONSTRAINT boards_country_id_title_unique UNIQUE (country_id, title);
  END IF;
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- grades
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS grades (
    id              BIGSERIAL PRIMARY KEY,
    grade_method_id BIGINT      NOT NULL REFERENCES grade_methods(id),
    title           VARCHAR(50) NOT NULL,
    age_range_start INTEGER     NULL,
    age_range_end   INTEGER     NULL,
    is_visible      BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS grades_grade_method_id_index ON grades(grade_method_id);
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'grades_grade_method_id_title_unique') THEN
    ALTER TABLE grades ADD CONSTRAINT grades_grade_method_id_title_unique UNIQUE (grade_method_id, title);
  END IF;
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- mediums
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS mediums (
    id          BIGSERIAL PRIMARY KEY,
    country_id  BIGINT       NOT NULL REFERENCES countries(id),
    board_id    BIGINT       NULL REFERENCES boards(id),
    title       VARCHAR(100) NOT NULL,
    language_id BIGINT       NULL REFERENCES languages(id),
    is_visible  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS mediums_country_id_index ON mediums(country_id);
CREATE INDEX IF NOT EXISTS mediums_board_id_index   ON mediums(board_id);
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'mediums_country_id_board_id_language_id_unique') THEN
    ALTER TABLE mediums ADD CONSTRAINT mediums_country_id_board_id_language_id_unique
        UNIQUE (country_id, board_id, language_id);
  END IF;
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- subjects
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS subjects (
    id           BIGSERIAL PRIMARY KEY,
    board_id     BIGINT       NOT NULL REFERENCES boards(id),
    medium_id    BIGINT       NOT NULL REFERENCES mediums(id),
    grade_id     BIGINT       NOT NULL REFERENCES grades(id),
    title        VARCHAR(100) NOT NULL,
    subject_code VARCHAR(50)  NULL,
    is_visible   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS subjects_board_id_grade_id_index ON subjects(board_id, grade_id);

-- ─────────────────────────────────────────────────────────────────────────────
-- users
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id                 BIGSERIAL PRIMARY KEY,
    uuid               UUID           NOT NULL DEFAULT uuid_generate_v4(),
    email              VARCHAR(255)   NULL,
    phone_number       VARCHAR(50)    NULL,
    password_hash      TEXT           NOT NULL,
    first_name         VARCHAR(100)   NULL,
    last_name          VARCHAR(100)   NULL,
    profile_photo      VARCHAR(500)   NULL,
    role               VARCHAR(50)    NOT NULL DEFAULT 'teacher',
    is_active          BOOLEAN        NOT NULL DEFAULT TRUE,
    is_verified        BOOLEAN        NULL,
    email_verified_at  TIMESTAMPTZ    NULL,
    last_login_at      TIMESTAMPTZ    NULL,
    created_at         TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ    NULL,
    preferable_subject VARCHAR(255)[] NOT NULL DEFAULT '{}',
    plateform_version  VARCHAR(255)   NOT NULL DEFAULT '0.0.0'
);

-- ─────────────────────────────────────────────────────────────────────────────
-- user_default
-- (current_state_id kept as plain BIGINT — states table was dropped in 000011)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_default (
    user_id            BIGINT      NOT NULL PRIMARY KEY REFERENCES users(id),
    current_country_id BIGINT      NULL REFERENCES countries(id),
    current_state_id   BIGINT      NULL,
    current_board_id   BIGINT      NULL REFERENCES boards(id),
    current_medium_id  BIGINT      NULL REFERENCES mediums(id),
    current_grade_id   BIGINT      NULL REFERENCES grades(id),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ─────────────────────────────────────────────────────────────────────────────
-- books
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS books (
    id               BIGSERIAL PRIMARY KEY,
    subject_id       BIGINT        NOT NULL REFERENCES subjects(id),
    title            VARCHAR(500)  NOT NULL,
    publisher        VARCHAR(255)  NULL,
    publication_year INTEGER       NULL,
    created_by       BIGINT        NULL REFERENCES users(id),
    is_public        BOOLEAN       NULL,
    is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
    file_path        VARCHAR(1000) NULL,
    total_pages      INTEGER       NULL,
    is_visible       BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS books_is_active_index ON books(is_active);

-- ─────────────────────────────────────────────────────────────────────────────
-- chapters
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS chapters (
    id              BIGSERIAL PRIMARY KEY,
    book_id         BIGINT        NOT NULL REFERENCES books(id),
    chapter_title   VARCHAR(500)  NOT NULL,
    content_summary TEXT          NULL,
    concept_tags    TEXT[]        NULL,
    embedding_id    VARCHAR(100)  NULL,
    is_visible      BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file_path       VARCHAR(1000) NULL
);
CREATE INDEX IF NOT EXISTS chapters_book_id_index      ON chapters(book_id);
CREATE INDEX IF NOT EXISTS chapters_concept_tags_index ON chapters(concept_tags);

-- ─────────────────────────────────────────────────────────────────────────────
-- generated_content
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS generated_content (
    id                          BIGSERIAL PRIMARY KEY,
    content_type                VARCHAR(50)   NOT NULL,
    chapter_id                  BIGINT        NOT NULL REFERENCES chapters(id),
    generated_by_user_id        BIGINT        NOT NULL REFERENCES users(id),
    generation_prompt           TEXT          NULL,
    generation_params           JSONB         NULL,
    generation_model            VARCHAR(50)   NULL,
    generation_duration_seconds INTEGER       NULL,
    file_url                    VARCHAR(1000) NOT NULL,
    file_path                   VARCHAR(1000) NULL,
    file_size_bytes             BIGINT        NULL,
    file_format                 VARCHAR(20)   NULL,
    thumbnail_url               VARCHAR(500)  NULL,
    title                       VARCHAR(255)  NULL,
    description                 TEXT          NULL,
    slide_count                 INTEGER       NULL,
    page_count                  INTEGER       NULL,
    question_count              INTEGER       NULL,
    embedding_id                VARCHAR(100)  NULL,
    embedding_generated_at      TIMESTAMPTZ   NULL,
    embedding_model             VARCHAR(50)   NULL,
    embedding_dimensions        INTEGER       NULL,
    concept_tags                TEXT[]        NULL,
    medium_id                   BIGINT        NOT NULL REFERENCES mediums(id),
    grade_id                    BIGINT        NOT NULL REFERENCES grades(id),
    subject_id                  BIGINT        NOT NULL REFERENCES subjects(id),
    board_id                    BIGINT        NOT NULL REFERENCES boards(id),
    state_id                    BIGINT        NULL,
    usage_count                 INTEGER       NOT NULL DEFAULT 0,
    download_count              INTEGER       NOT NULL DEFAULT 0,
    view_count                  INTEGER       NOT NULL DEFAULT 0,
    recommendation_accept_count INTEGER       NOT NULL DEFAULT 0,
    recommendation_reject_count INTEGER       NOT NULL DEFAULT 0,
    quality_score               DECIMAL(3,2)  NULL,
    average_rating              DECIMAL(3,2)  NULL,
    rating_count                INTEGER       NOT NULL DEFAULT 0,
    is_reusable                 BOOLEAN       NOT NULL DEFAULT TRUE,
    is_anonymous                BOOLEAN       NOT NULL DEFAULT TRUE,
    is_public                   BOOLEAN       NOT NULL DEFAULT TRUE,
    share_scope                 VARCHAR(50)   NOT NULL DEFAULT 'global',
    status                      VARCHAR(50)   NOT NULL DEFAULT 'active',
    created_at                  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at                TIMESTAMPTZ   NULL
);
CREATE INDEX IF NOT EXISTS generated_content_content_type_index              ON generated_content(content_type);
CREATE INDEX IF NOT EXISTS generated_content_chapter_id_index                ON generated_content(chapter_id);
CREATE INDEX IF NOT EXISTS generated_content_generated_by_user_id_index      ON generated_content(generated_by_user_id);
CREATE INDEX IF NOT EXISTS generated_content_medium_id_grade_id_subject_id_index
    ON generated_content(medium_id, grade_id, subject_id);
CREATE INDEX IF NOT EXISTS generated_content_quality_score_usage_count_index ON generated_content(quality_score, usage_count);
CREATE INDEX IF NOT EXISTS generated_content_concept_tags_index              ON generated_content(concept_tags);
CREATE INDEX IF NOT EXISTS generated_content_state_id_index                  ON generated_content(state_id);

-- =============================================================================
-- Done. All tables created.
-- =============================================================================
