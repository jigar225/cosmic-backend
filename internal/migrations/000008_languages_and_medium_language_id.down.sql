-- Restore previous uniqueness and drop language linkage.

-- Re-introduce the old unique constraint on (country_id, board_id, title) if missing.
ALTER TABLE mediums
    ADD CONSTRAINT IF NOT EXISTS mediums_country_id_board_id_title_unique
        UNIQUE (country_id, board_id, title);

-- Drop the new unique constraint on (country_id, board_id, language_id).
ALTER TABLE mediums
    DROP CONSTRAINT IF EXISTS mediums_country_id_board_id_language_id_unique;

-- Drop the language_id column from mediums.
ALTER TABLE mediums
    DROP COLUMN IF EXISTS language_id;

-- Drop the global languages table.
DROP TABLE IF EXISTS languages;

