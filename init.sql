create extension if not exists citext;

drop table if exists forums cascade;
drop table if exists users cascade;
drop table if exists threads cascade;

create table users
(
    id       bigserial not null,
    nickname citext    not null primary key,
    fullname varchar   not null,
    email    citext    not null unique,
    about    text
);

-- select *
-- from pg_indexes
-- where tablename = 'forums';

create table forums
(
    id      bigserial not null,
    posts   int       not null default 0,
    slug    citext    not null primary key,
    threads int       not null default 0,
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
    votes      int         default 0,
    created_at timestamptz default now()
);
