create extension if not exists citext;

drop table if exists forums;
drop table if exists users;

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
    usr     citext    not null references users(nickname)
);
