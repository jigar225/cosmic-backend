-- Allow file_path to be NULL for "unit" books (chapters uploaded separately; no whole-book file).
-- When user uploads a full book, we set file_path.
ALTER TABLE books ALTER COLUMN file_path DROP NOT NULL;
