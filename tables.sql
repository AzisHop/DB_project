ALTER USER postgres WITH ENCRYPTED PASSWORD 'admin';
create extension if not exists citext;
drop table if exists votes;
drop table if exists post;
drop table if exists thread;
drop table if exists allUsersForum;
drop table if exists forum;
drop table if exists userForum;

DROP INDEX if exists Users;
DROP INDEX if exists threadSlug;
DROP INDEX if exists threadCreated;
DROP INDEX if exists threadForumHash;
DROP INDEX if exists threads;
DROP INDEX if exists postParent;
DROP INDEX if exists postPath;
DROP INDEX if exists voteIndex;
DROP INDEX if exists usersAll;
DROP INDEX if exists nickname;
DROP INDEX if exists postPath2;

CREATE TABLE userForum
(
    nickname CITEXT collate "POSIX" PRIMARY KEY NOT NULL,
    fullname TEXT                               NOT NULL,
    about    TEXT,
    email    CITEXT             UNIQUE          NOT NULL
);

CREATE INDEX  Users ON userforum (nickname, fullname, about, email);
CREATE INDEX if not exists nickname on userforum using hash (nickname);

CREATE TABLE forum
(
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT          NOT NULL,
    "user"  citext        NOT NULL,
    slug    citext UNIQUE NOT NULL,
    posts   BIGINT        NOT NULL DEFAULT 0,
    threads BIGINT        NOT NULL DEFAULT 0,
    FOREIGN KEY ("user") REFERENCES userForum (nickname)
);

CREATE TABLE thread
(
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT                     NOT NULL,
    author  citext                   NOT NULL,
    forum   citext                   NOT NULL,
    message TEXT                     NOT NULL,
    votes   BIGINT                   NOT NULL DEFAULT 0,
    slug    citext UNIQUE,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (author)
        REFERENCES userForum (nickname),
    FOREIGN KEY (forum)
        REFERENCES forum (slug)
);


CREATE OR REPLACE FUNCTION postTriggerFunc()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum SET threads = threads + 1 WHERE slug = NEW.forum;

    INSERT INTO allUsersForum(forum, nickname, fullname, about, email)
    SELECT NEW.forum, nickname, fullname, about, email
    FROM userForum WHERE nickname = NEW.author
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS forum_thread ON thread;
CREATE TRIGGER forum_thread
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE postTriggerFunc();

CREATE INDEX  threadSlug ON thread using hash (slug);
CREATE INDEX  threadCreated ON thread (created);
CREATE INDEX  threadForumHash ON thread using hash (forum);
CREATE INDEX  threads on thread (forum, slug, created,title, author, message, votes);


CREATE UNLOGGED TABLE post
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
    FOREIGN KEY (author) REFERENCES userForum (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    FOREIGN KEY (thread) REFERENCES thread (id)
);

CREATE OR REPLACE FUNCTION postTriggerFunction()
    RETURNS TRIGGER AS
$$
DECLARE
    pP BIGINT[];
BEGIN
    IF (NEW.parent IS NULL)
        THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM post
        WHERE id = new.parent INTO pP;
        NEW.path := NEW.path
                        || pP
                        || new.id;
    END IF;

    UPDATE forum SET posts = posts + 1 WHERE slug = NEW.forum;

    INSERT INTO allUsersForum(forum, nickname, fullname, about, email)
    SELECT NEW.forum, nickname, fullname, about, email
    FROM userForum
    WHERE nickname = NEW.author
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS postTrigger ON post;
CREATE TRIGGER postTrigger
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE postTriggerFunction();

CREATE INDEX  postParent on post (thread, parent);
CREATE INDEX  postPath on post ((path[1]), id);
CREATE INDEX postPath2 on post using gin (path);

drop table if exists allUsersForum;
CREATE UNLOGGED TABLE allUsersForum
(
    forum    citext                 NOT NULL,
    nickname citext collate "POSIX" NOT NULL,
    fullname TEXT                   NOT NULL,
    about    TEXT,
    email    citext                 NOT NULL,
    FOREIGN KEY (forum)
        REFERENCES forum (slug),
    FOREIGN KEY (nickname)
        REFERENCES userForum (nickname),
    PRIMARY KEY (nickname, forum)
);

CREATE INDEX usersAll on allUsersForum (forum, nickname, fullname, about, email)

drop table if exists votes;
CREATE UNLOGGED TABLE votes (
                                thread   bigint    NOT NULL,
                                nickname citext    NOT NULL,
                                voice    BIGINT    NOT NULL,
                                FOREIGN KEY (thread)
                                    REFERENCES thread (id),
                                FOREIGN KEY (nickname)
                                    REFERENCES userForum (nickname),
                                PRIMARY KEY (thread, nickname)
);



-- CREATE INDEX  Users ON userforum (nickname, fullname, about, email);
-- CREATE INDEX if not exists nickname on userforum using hash (nickname);
-- CREATE INDEX  postParent on post (thread, parent);
-- CREATE INDEX  postPath on post ((path[1]), id);
-- CREATE INDEX postPath2 on post using gin (path);
-- CREATE INDEX  threadSlug ON thread using hash (slug);
-- CREATE INDEX  threadCreated ON thread (created);
-- CREATE INDEX  threadForumHash ON thread using hash (forum);
-- CREATE INDEX  threads on thread (forum, slug, created,title, author, message, votes);
CREATE INDEX voteIndex on votes (thread, nickname, voice);
-- CREATE INDEX usersAll on allUsersForum (forum, nickname, fullname, about, email)





