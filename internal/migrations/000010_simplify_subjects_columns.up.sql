ALTER TABLE subjects
    DROP COLUMN IF EXISTS subject_code,
    DROP COLUMN IF EXISTS sequence_order,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS created_by;

