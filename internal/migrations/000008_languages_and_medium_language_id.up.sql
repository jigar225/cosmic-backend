CREATE TABLE languages (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    is_visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE mediums
    ADD COLUMN language_id BIGINT NULL REFERENCES languages(id);

-- Backfill languages from existing mediums using language_code where present.
INSERT INTO languages (code, name)
SELECT DISTINCT m.language_code, m.title
FROM mediums m
WHERE m.language_code IS NOT NULL
ON CONFLICT (code) DO NOTHING;

-- Populate medium.language_id based on the backfilled languages.
UPDATE mediums m
SET language_id = l.id
FROM languages l
WHERE m.language_code IS NOT NULL
  AND m.language_code = l.code;

-- Ensure one medium per (country, board, language).
ALTER TABLE mediums
    ADD CONSTRAINT mediums_country_id_board_id_language_id_unique
        UNIQUE (country_id, board_id, language_id);

-- Old unique by title is no longer needed for new design.
ALTER TABLE mediums
    DROP CONSTRAINT IF EXISTS mediums_country_id_board_id_title_unique;

