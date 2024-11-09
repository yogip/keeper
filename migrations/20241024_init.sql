-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id       BIGSERIAL PRIMARY KEY,
    email    VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL
);

-- Folders table to group secrets together
CREATE TABLE IF NOT EXISTS folders(
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT REFERENCES users(id),
    name    VARCHAR,
    UNIQUE (id, user_id)
);

-- Secrets tables
CREATE TABLE IF NOT EXISTS passwords(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id),
    folder_id  BIGINT REFERENCES folders(id) DEFAULT NULL,
    name       VARCHAR NOT NULL,
    login      VARCHAR NOT NULL,
    password   VARCHAR NOT NULL,
    sc_version INT NOT NULL,
    sc         VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS notes(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id),
    folder_id  BIGINT REFERENCES folders(id) DEFAULT NULL,
    name       VARCHAR NOT NULL,
    note       text NOT NULL,
    sc_version INT NOT NULL,
    sc         VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS files(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id),
    folder_id  BIGINT REFERENCES folders(id) DEFAULT NULL,
    name       VARCHAR NOT NULL,
    path       VARCHAR NOT NULL,
    sc_version INT NOT NULL,
    sc         VARCHAR NOT NULL
);


-- Metadata table to store tags
CREATE TYPE TAG_RELATIONS AS ENUM ('password', 'note', 'file');

CREATE TABLE IF NOT EXISTS tags(
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT REFERENCES users(id),
    relation TAG_RELATIONS NOT NULL,
    name     VARCHAR NOT NULL,
    value    VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS passwords;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS folders;
DROP TYPE IF EXISTS TAG_RELATIONS;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
