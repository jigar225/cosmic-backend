-- Replace signup_method (single) with signup_methods (ordered array).
-- Order = display order: first = priority (top), rest below.

ALTER TABLE countries ADD COLUMN IF NOT EXISTS signup_methods TEXT[] NULL;

UPDATE countries
SET signup_methods = CASE
  WHEN signup_method IS NOT NULL AND signup_method != '' THEN ARRAY[signup_method]
  ELSE ARRAY['email']::TEXT[]
END
WHERE signup_methods IS NULL;

ALTER TABLE countries ALTER COLUMN signup_methods SET DEFAULT ARRAY['email']::TEXT[];
ALTER TABLE countries ALTER COLUMN signup_methods SET NOT NULL;

ALTER TABLE countries DROP COLUMN IF EXISTS signup_method;
