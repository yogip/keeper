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
    note         VARCHAR NOT NULL default '',
    secret_type  SECRET_TYPE NOT NULL,
    sc_version   INT NOT NULL,
    sc           VARCHAR NOT NULL
);
CREATE INDEX spgist on secrets using spgist (name text_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS SECRET_TYPE;
-- +goose StatementEnd
