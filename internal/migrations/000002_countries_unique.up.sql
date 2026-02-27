-- Ensure country_code is unique per country.
ALTER TABLE countries
    ADD CONSTRAINT countries_country_code_unique UNIQUE (country_code);

