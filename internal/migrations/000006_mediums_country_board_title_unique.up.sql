ALTER TABLE mediums
    ADD CONSTRAINT mediums_country_id_board_id_title_unique
        UNIQUE (country_id, board_id, title);

