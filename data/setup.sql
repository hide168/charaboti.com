drop table users;
drop table sessions;
drop table characters;

create table users (
    id          serial primary key,
    uuid        varchar(64) not null unique,
    name        varchar(255),
    email       varchar(255) not null unique,
    password    varchar(255) not null,
    icon        varchar(255) not null,
    created_at  timestamp not null
);

create table sessions (
    id          serial primary key,
    uuid        varchar(64) not null unique,
    email       varchar(255),
    user_id     integer references users(id),
    created_at  timestamp not null
);

create table characters (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  text       text,
  user_id    integer references users(id),
  image      varchar(255) not null,
  created_at timestamp not null       
);