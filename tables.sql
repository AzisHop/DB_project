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

drop table if exists post;
CREATE TABLE post
(
    id       BIGSERIAL PRIMARY KEY,
    parent   BIGINT                   NOT NULL,
    author   text                     NOT NULL,
    message  TEXT                     NOT NULL,
    isEdited BOOLEAN                  NOT NULL DEFAULT false,
    forum    text                     NOT NULL,
    thread   BIGINT                   NOT NULL,
    created  TIMESTAMP WITH TIME ZONE NOT NULL,
    path     BIGINT[]                 NOT NULL DEFAULT ARRAY []::INTEGER[],
    FOREIGN KEY (author)
        REFERENCES userForum (nickname),
    FOREIGN KEY (forum)
        REFERENCES forum (slug),
    FOREIGN KEY (thread)
        REFERENCES thread (id)
);
drop table if exists allUsersForum;
CREATE UNLOGGED TABLE allUsersForum
(
    forum    text                 NOT NULL,
    nickname text               collate "POSIX"     NOT NULL,
    fullname TEXT                   NOT NULL,
    about    TEXT,
    email    text                 NOT NULL,
    FOREIGN KEY (forum)
        REFERENCES forum (slug),
    FOREIGN KEY (nickname)
        REFERENCES userForum (nickname),
    PRIMARY KEY (nickname, forum)
);

SELECT nickname, fullname, about, email
FROM allUsersForum
WHERE forum = 'pirate-stories'

SELECT DISTINCT nickname, fullname, about, email
FROM userForum uF
INNER JOIN  thread t on uF.nickname = t.author AND t.forum = 'pirate-stories'
INNER JOIN post p on uF.nickname = p.author AND p.forum = 'pirate-stories'
-- INNER JOIN forum f on uF.nickname = f."user" AND f.slug = 'pirate-stories'




SELECT id, title, author, forum, message, votes, slug, created
FROM thread tr
WHERE forum = 'pirate-stories' AND created = '1918-06-13 18:54:05.031000'
ORDER BY created DESC LIMIT 100

SELECT slug FROM thread WHERE id = 3