ALTER TABLE subjects
    ADD COLUMN subject_code VARCHAR(50) NULL,
    ADD COLUMN sequence_order INTEGER NULL,
    ADD COLUMN description TEXT NULL,
    ADD COLUMN created_by BIGINT NULL;

