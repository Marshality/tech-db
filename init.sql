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
    posts   int       not null default 0,
    slug    citext    not null primary key,
    threads int       not null default 0,
    title   varchar   not null,
    usr     citext    not null references users (nickname)
);

create table threads
(
    id         bigserial not null primary key,
    slug       citext not null,
    author     citext    not null references users (nickname),
    forum      citext    not null references forums (slug),
    message    text,
    title      varchar   not null,
    votes      int         default 0,
    created_at timestamptz default now()
);

drop function if exists threadsCounter;
create or replace function threadsCounter()
    returns trigger as
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

drop function if exists postsCounter;
create or replace function postsCounter()
    returns trigger as
$$
begin
    update forums
    set posts = posts + 1
    where slug = new.forum;

    return null;
end;
$$ language plpgsql;

drop trigger if exists postsIncrementer on threads;
create trigger postsIncrementer
    after insert
    on threads
    for each row
execute procedure postsCounter();

create table vote
(
    id bigserial primary key,
    thread int not null references threads(id),
    nickname citext not null references users(nickname),
    voice int not null
);

create unique index on vote (thread, nickname);

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
