-- Restore single signup_method column (first element of array).
ALTER TABLE countries ADD COLUMN IF NOT EXISTS signup_method VARCHAR(50) NOT NULL DEFAULT 'email';
UPDATE countries SET signup_method = COALESCE(signup_methods[1], 'email') WHERE array_length(signup_methods, 1) > 0;
ALTER TABLE countries DROP COLUMN IF EXISTS signup_methods;
