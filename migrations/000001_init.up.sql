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
    user_id integer not null,
    list_id integer not null,
    is_admin bool default false,
    CONSTRAINT fk_users_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_lists_id FOREIGN KEY(list_id) REFERENCES lists(id) ON DELETE CASCADE
);

create table items (
    id serial primary key,
    list_id integer not null,
    title varchar(255) not null,
    description text not null,
    done bool default false,
    CONSTRAINT fk_lists_id FOREIGN KEY(list_id) REFERENCES lists(id) ON DELETE CASCADE
);