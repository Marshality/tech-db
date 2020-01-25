create extension if not exists citext;

drop table if exists forums cascade;
drop table if exists users cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists vote cascade;

create unlogged table users
(
    id       bigserial not null,
    nickname citext    not null primary key,
    fullname varchar   not null,
    email    citext    not null unique,
    about    text
);

create unlogged table forums
(
    id      bigserial not null,
    posts   integer   not null default 0,
    slug    citext    not null primary key,
    threads integer   not null default 0,
    title   varchar   not null,
    usr     citext    not null references users (nickname)
);

create unlogged table threads
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

create index on user_forum (forum_slug);

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

create unlogged table posts
(
    id         integer   not null primary key,
    forum      citext    not null references forums (slug),
    thread     bigint    not null references threads (id),
    author     citext    not null references users (nickname),
    message    text,
    parent     integer   not null default 0,
    is_edited  bool      not null default false,
    created_at timestamptz        default now(),
    path       bigint array
);

CREATE SEQUENCE post_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE post_id_seq OWNED BY posts.id;
SELECT pg_catalog.setval('post_id_seq', 1, false);

create index on posts (thread);
create index ON posts using gin (path);

create unlogged table vote
(
    id       bigserial primary key,
    thread   integer not null references threads (id),
    nickname citext  not null references users (nickname),
    voice    integer not null,
    constraint unique_vote unique (nickname, thread)
);

-- create index on vote (thread, nickname);

select id, path
from posts;

