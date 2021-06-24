ALTER USER postgres WITH ENCRYPTED PASSWORD 'admin';
create extension if not exists citext;
drop table if exists votes;
drop table if exists post;
drop table if exists thread;
drop table if exists allUsersForum;
drop table if exists forum;
drop table if exists userForum;


CREATE TABLE userForum
(
    id       BIGSERIAL PRIMARY KEY,
    nickname citext UNIQUE NOT NULL,
    fullname TEXT        NOT NULL,
    about    TEXT        NOT NULL,
    email    citext UNIQUE NOT NULL
);
drop table if exists forum;
CREATE TABLE forum
(
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT   NOT NULL,
    "user"  citext NOT NULL,
    slug    citext UNIQUE NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    FOREIGN KEY ("user") REFERENCES userForum (nickname)
);

drop table if exists thread;
CREATE TABLE thread
(
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT                     NOT NULL,
    author  citext                   NOT NULL,
    forum   citext                   NOT NULL,
    message TEXT                     NOT NULL,
    votes   BIGINT                   NOT NULL DEFAULT 0,
    slug    citext UNIQUE DEFAULT NULL,
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
    author   citext                     NOT NULL,
    message  TEXT                     NOT NULL,
    isEdited BOOLEAN                  NOT NULL DEFAULT false,
    forum    citext                     NOT NULL,
    thread   BIGINT                   NOT NULL,
    created  TIMESTAMP WITH TIME ZONE NOT NULL,
    path BIGINT[] DEFAULT '{}',
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
    forum    citext                 NOT NULL,
    nickname citext               collate "POSIX"     NOT NULL,
    fullname TEXT                   NOT NULL,
    about    TEXT,
    email    citext                 NOT NULL,
    FOREIGN KEY (forum)
        REFERENCES forum (slug),
    FOREIGN KEY (nickname)
        REFERENCES userForum (nickname),
    PRIMARY KEY (nickname, forum)
);
drop table if exists votes;
CREATE TABLE votes (
    thread INT NOT NULL,
    voice INT NOT NULL,
    nickname citext NOT NULL,
    FOREIGN KEY (thread) REFERENCES thread (id),
    FOREIGN KEY (nickname) REFERENCES userForum(nickname),
    UNIQUE (thread, nickname)
);

CREATE OR REPLACE FUNCTION threadTriggerFunc()
    RETURNS trigger AS
$$
BEGIN
    INSERT INTO allUsersForum(forum, nickname,fullname, about, email)
        SELECT new.forum, nickname, fullname, about, email
        FROM userForum
            WHERE nickname = new.author
    ON CONFLICT DO NOTHING;

    UPDATE forum SET threads = threads + 1
    WHERE slug = new.forum;

    RETURN NEW;
END;
$$
    LANGUAGE 'plpgsql';

drop trigger IF EXISTS threadTrigger on thread;

CREATE TRIGGER threadTrigger
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE threadTriggerFunc();

CREATE OR REPLACE FUNCTION postTriggerFunc()
    RETURNS trigger AS
$$
BEGIN
    INSERT INTO allUsersForum(forum, nickname,fullname, about, email)
        SELECT new.forum, nickname, fullname, about, email
        FROM userForum
            WHERE nickname = new.author
    ON CONFLICT DO NOTHING;

    UPDATE forum SET posts = forum.posts + 1
    WHERE slug = new.forum;

    NEW.path = (SELECT path FROM post WHERE id = NEW.parent LIMIT 1) || NEW.id;

    RETURN NEW;
END;
$$
    LANGUAGE 'plpgsql';

drop trigger IF EXISTS postTrigger on post;

CREATE TRIGGER postTrigger
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE postTriggerFunc();

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

SELECT id, title, author, forum, message, votes, coalesce(slug,''),created
FROM thread tr
WHERE forum = 'KH0Nm-SP_L5J8'
  AND date_part('year', created)  = date_part('year', date '2021-08-31T12:27:00.307Z') ORDER BY created DESC LIMIT 4

SELECT id, title, author, forum, message, votes, coalesce(slug,''),created FROM thread tr WHERE forum = 'KH0Nm-SP_L5J8' AND date_part('year', created)  = date_part('year', date '2021-08-31T12:27:00.307Z') ORDER BY created DESC LIMIT 4



SELECT nickname FROM userForum WHERE nickname = 'mira.rP28QfR9DIfFJU';

SELECT nickname FROM userForum WHERE nickname = 'aut.s1y035NL1c5Cj';

SELECT title, "user", coalesce(slug, ''), posts, threads FROM forum WHERE slug = 'cC32xovt-jJfK';

SELECT slug FROM forum WHERE slug = 'sOjqtEFg-FCF86';


SELECT nickname FROM userForum WHERE nickname = 'o.doy3Qyi05C55RU'

DROP INDEX if exists Users;
DROP INDEX if exists threadSlug;
DROP INDEX if exists threadCreated;
DROP INDEX if exists threadForumHash;
DROP INDEX if exists threads;
DROP INDEX if exists postParent;
DROP INDEX if exists postPath;
DROP INDEX if exists voteIndex;
DROP INDEX if exists usersAll;

CREATE INDEX  Users ON userforum (nickname, fullname, about, email);
create INDEX  postParent on post (thread, parent);
create INDEX  postPath on post ((path[1]), id);
CREATE INDEX  threadSlug ON thread using hash (slug);
CREATE INDEX  threadCreated ON thread (created);
CREATE INDEX  threadForumHash ON thread using hash (forum);
CREATE INDEX  threads on thread (forum, slug, created,title, author, message, votes);
CREATE INDEX voteIndex on votes (thread, nickname, voice);
CREATE INDEX usersAll on allUsersForum (forum, nickname, fullname, about, email)
