DROP TABLE IF EXISTS blocks;
CREATE TABLE IF NOT EXISTS blocks (
    network text,
    geoname_id integer,
    registered_country_geoname_id integer,
    represented_country_geoname_id integer,
    is_anonymous_proxy integer,
    is_satellite_provider integer,
    postal_code text,
    latitude real,
    longitude real,
    accuracy_radius integer
)