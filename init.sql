create extension if not exists citext;

drop table if exists forums cascade;
drop table if exists users cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists vote cascade;

create table users
(
    id       bigserial not null,
    nickname citext    not null primary key,
    fullname varchar   not null,
    email    citext    not null unique,
    about    text
);

create table forums
(
    id      bigserial not null,
    posts   integer   not null default 0,
    slug    citext    not null primary key,
    threads integer   not null default 0,
    title   varchar   not null,
    usr     citext    not null references users (nickname)
);

create table threads
(
    id         bigserial not null primary key,
    slug       citext    not null,
    author     citext    not null references users (nickname),
    forum      citext    not null references forums (slug),
    message    text,
    title      varchar   not null,
    votes      integer     default 0,
    created_at timestamptz default now()
);

create index on threads (slug);
create index on threads (created_at, forum);

drop table if exists user_forum cascade;
create table user_forum
(
    user_id    integer,
    forum_slug citext,
    primary key (forum_slug, user_id)
);

create index index_on_forum_user on user_forum (forum_slug);

drop function if exists threadsCounter;
create or replace function threadsCounter()
    returns trigger as
$$
begin
    insert into user_forum (forum_slug, user_id)
    values (new.forum, (select id from users where nickname = new.author))
    on conflict do nothing;

    update forums
    set threads = threads + 1
    where slug = new.forum;

    return null;
end;
$$ language plpgsql;

drop trigger if exists threadsIncrementer on threads;
create trigger threadsIncrementer
    after insert
    on threads
    for each row
execute procedure threadsCounter();

create table posts
(
    id         bigserial not null primary key,
    forum      citext    not null references forums (slug),
    thread     bigint    not null references threads (id),
    author     citext    not null references users (nickname),
    message    text,
    parent     integer   not null default 0,
    is_edited  bool      not null default false,
    created_at timestamptz        default now(),
    path       bigint array
);

create index on posts (thread);

drop function if exists postsCounter;
create or replace function postsCounter()
    returns trigger as
$$
begin
    insert into user_forum (forum_slug, user_id)
    values (new.forum, (select id from users where nickname = new.author))
    on conflict do nothing;

    update forums
    set posts = posts + 1
    where slug = new.forum;

    return null;
end;
$$ language plpgsql;

drop trigger if exists postsIncrementer on threads;
create trigger postsIncrementer
    after insert
    on posts
    for each row
execute procedure postsCounter();

-- post path
drop function if exists setPath;
create or replace function setPath()
    returns trigger as
$$
begin
    if new.parent = 0 THEN
        UPDATE posts
        SET path = ARRAY [new.id]
        WHERE id = new.id;
    ELSE
        UPDATE posts
        SET path = array_append((SELECT path FROM posts WHERE id = new.parent), new.id)
        WHERE id = new.id;
    END IF;
    RETURN new;
END;
$$
    LANGUAGE plpgsql;

drop trigger if exists pathSetter on posts;
create trigger pathSetter
    after insert
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE setPath();
-- post path

create table vote
(
    id       bigserial primary key,
    thread   integer not null references threads (id),
    nickname citext  not null references users (nickname),
    voice    integer not null,
    constraint unique_vote unique (nickname, thread)
);

create index on vote (thread, nickname);

drop function if exists voteInsert cascade;
create or replace function voteInsert()
    returns trigger as
$$
begin
    update threads
    set votes = votes + new.voice
    where id = new.thread;

    return null;
end;
$$ language plpgsql;

drop trigger if exists voteInserter on threads cascade;
create trigger voteInserter
    after insert
    on vote
    for each row
execute procedure voteInsert();

drop function if exists voteUpdate cascade;
create or replace function voteUpdate()
    returns trigger as
$$
begin
    update threads
    set votes = votes - old.voice + new.voice
    where id = new.thread;

    return null;
end;
$$ language plpgsql;

drop trigger if exists voteUpdater on threads cascade;
create trigger voteUpdater
    after update
    on vote
    for each row
execute procedure voteUpdate();

select id, path
from posts;
