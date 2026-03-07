-- Add display_order to chapters for stable front-end ordering.
-- Backfill existing rows using their row number per book (ordered by id).
ALTER TABLE chapters ADD COLUMN IF NOT EXISTS display_order INTEGER NOT NULL DEFAULT 0;

UPDATE chapters c
SET display_order = sub.rn
FROM (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY book_id ORDER BY id) AS rn
    FROM chapters
) sub
WHERE c.id = sub.id;
