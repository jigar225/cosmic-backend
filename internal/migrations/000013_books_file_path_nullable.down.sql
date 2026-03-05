-- Restore NOT NULL: set any NULL file_path to empty string first.
UPDATE books SET file_path = '' WHERE file_path IS NULL;
ALTER TABLE books ALTER COLUMN file_path SET NOT NULL;
