create extension if not exists citext;

drop table if exists forums cascade;
drop table if exists users cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;

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
    slug       citext unique,
    author     citext    not null references users (nickname),
    forum      citext    not null references forums (slug),
    message    text,
    title      varchar   not null,
    votes      int         default 0,
    created_at timestamptz default now()
);

drop function if exists threadsCounter;
create or replace function threadsCounter()
    returns trigger AS
$$
begin
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
    message    text      not null,
    parent     int       not null default 0,
    is_edited  bool      not null default false,
    created_at timestamptz        default now()
);
