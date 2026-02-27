-- One board title per country (e.g. only one "CBSE" per country).
ALTER TABLE boards
    ADD CONSTRAINT boards_country_id_title_unique UNIQUE (country_id, title);
