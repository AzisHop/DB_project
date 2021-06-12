drop table if exists userForum;

CREATE TABLE userForum
(
    id       BIGSERIAL PRIMARY KEY,
    nickname TEXT UNIQUE NOT NULL,
    fullname TEXT        NOT NULL,
    about    TEXT        NOT NULL,
    email    TEXT UNIQUE NOT NULL
)