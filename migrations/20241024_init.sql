-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id       BIGSERIAL PRIMARY KEY,
    email    VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL
);

CREATE TYPE SECRET_TYPE AS ENUM ('password', 'note', 'file', 'card');

-- Secrets tables
CREATE TABLE IF NOT EXISTS secrets(
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT REFERENCES users(id),
    name         VARCHAR NOT NULL,
    payload      VARCHAR NOT NULL,
    secret_type  SECRET_TYPE NOT NULL,
    sc_version   INT NOT NULL,
    sc           VARCHAR NOT NULL
);
CREATE INDEX spgist on secrets using spgist (name text_ops);

-- Metadata table to store tags
CREATE TABLE IF NOT EXISTS tags(
    id           BIGSERIAL PRIMARY KEY,
    secret_id    BIGINT REFERENCES secrets(id),
    user_id      BIGINT REFERENCES users(id),
    name         VARCHAR NOT NULL,
    value        VARCHAR NOT NULL
);

-- CREATE TABLE IF NOT EXISTS passwords(
--     id         BIGSERIAL PRIMARY KEY,
--     user_id    BIGINT REFERENCES users(id),
--     name       VARCHAR NOT NULL,
--     login      VARCHAR NOT NULL,
--     password   VARCHAR NOT NULL,
--     sc_version INT NOT NULL,
--     sc         VARCHAR NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS notes(
--     id         BIGSERIAL PRIMARY KEY,
--     user_id    BIGINT REFERENCES users(id),
--     name       VARCHAR NOT NULL,
--     note       text NOT NULL,
--     sc_version INT NOT NULL,
--     sc         VARCHAR NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS files(
--     id         BIGSERIAL PRIMARY KEY,
--     user_id    BIGINT REFERENCES users(id),
--     name       VARCHAR NOT NULL,
--     path       VARCHAR NOT NULL,
--     sc_version INT NOT NULL,
--     sc         VARCHAR NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS cards(
--     id          BIGSERIAL PRIMARY KEY,
--     user_id     BIGINT REFERENCES users(id),
--     name        VARCHAR NOT NULL,
--     payload     VARCHAR NOT NULL,
--     sc_version  INT NOT NULL,
--     sc          VARCHAR NOT NULL
-- );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS SECRET_TYPE;
-- +goose StatementEnd
