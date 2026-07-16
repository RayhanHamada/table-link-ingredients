-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied.

CREATE TABLE IF NOT EXISTS tm_ingredient (
    uuid         uuid         NOT NULL,
    name         varchar(255) NOT NULL,
    cause_alergy bool         NOT NULL,
    type         int4         NOT NULL,
    status       int4         DEFAULT 0,
    created_at   timestamp(6),
    updated_at   timestamp(6),
    deleted_at   timestamp(6),
    CONSTRAINT tm_ingredient_pkey PRIMARY KEY (uuid)
);

COMMENT ON COLUMN tm_ingredient.type   IS '0 (none), 1 (veggie), 2 (vegan)';
COMMENT ON COLUMN tm_ingredient.status IS '0 (inactive), 1 (active)';

CREATE TABLE IF NOT EXISTS tm_item (
    uuid       uuid          NOT NULL,
    name       varchar(255)  NOT NULL,
    price      numeric(10,2) NOT NULL,
    status     int4          DEFAULT 0,
    created_at timestamp(6),
    updated_at timestamp(6),
    deleted_at timestamp(6),
    CONSTRAINT tm_item_pkey PRIMARY KEY (uuid)
);

CREATE TABLE IF NOT EXISTS tm_item_ingredient (
    uuid_item       uuid NOT NULL,
    uuid_ingredient uuid NOT NULL,
    CONSTRAINT tm_item_ingredient_pkey PRIMARY KEY (uuid_item, uuid_ingredient),
    CONSTRAINT fk_item_ingredient_item       FOREIGN KEY (uuid_item)       REFERENCES tm_item (uuid),
    CONSTRAINT fk_item_ingredient_ingredient FOREIGN KEY (uuid_ingredient) REFERENCES tm_ingredient (uuid)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back.

DROP TABLE IF EXISTS tm_item_ingredient;
DROP TABLE IF EXISTS tm_item;
DROP TABLE IF EXISTS tm_ingredient;
