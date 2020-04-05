DROP TABLE IF EXISTS locations;
CREATE TABLE IF NOT EXISTS locations (
    geoname_id integer,
    locale_code text,
    continent_code text,
    continent_name text,
    country_iso_code text,
    country_name text,
    subdivision_1_iso_code text,
    subdivision_1_name text,
    subdivision_2_iso_code text,
    subdivision_2_name text,
    city_name text,
    metro_code text,
    time_zone text,
    is_in_european_union integer
)