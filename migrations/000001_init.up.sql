create table users (
    id serial primary key,
    email varchar(256) not null unique,
    password_hash varchar(256) not null,
    is_activated bool default false not null,
    activated_link varchar(128) not null
);

create table lists (
    id serial primary key,
    title varchar(255) not null,
    description text not null
);

create table users_lists (
    user_id integer references users(id),
    list_id integer references lists(id),
    is_admin bool default false
);

create table items (
    id serial primary key,
    list_id integer references lists(id),
    title varchar(255) not null,
    description text not null,
    done bool default false
);