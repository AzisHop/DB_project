drop table if exists userForum;


CREATE TABLE userForum
(
    id       BIGSERIAL PRIMARY KEY,
    nickname TEXT UNIQUE NOT NULL,
    fullname TEXT        NOT NULL,
    about    TEXT        NOT NULL,
    email    TEXT UNIQUE NOT NULL
);
drop table if exists forum;
CREATE TABLE forum
(
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT   NOT NULL,
    "user"  text NOT NULL,
    slug    text UNIQUE NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    FOREIGN KEY ("user") REFERENCES userForum (nickname)
);

drop table if exists thread;
CREATE TABLE thread
(
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT                     NOT NULL,
    author  text                   NOT NULL,
    forum   text                   NOT NULL,
    message TEXT                     NOT NULL,
    votes   BIGINT                   NOT NULL DEFAULT 0,
    slug    text UNIQUE DEFAULT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (author)
        REFERENCES userForum (nickname),
    FOREIGN KEY (forum)
        REFERENCES forum (slug)
);
