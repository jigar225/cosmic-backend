ALTER TABLE grades
    ADD CONSTRAINT grades_grade_method_id_title_unique UNIQUE (grade_method_id, title);

