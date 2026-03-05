-- Add PDF storage (S3 key) for chapters.
ALTER TABLE chapters ADD COLUMN IF NOT EXISTS file_path VARCHAR(1000) NULL;
