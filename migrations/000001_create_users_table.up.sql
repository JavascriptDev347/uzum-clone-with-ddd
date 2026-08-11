create type user_role as enum ('customer', 'admin');

create table users(
    id uuid primary key,
    email varchar(255) not null unique,
    password_hash varchar(255) not null,
    role user_role not null default 'customer',
    created_at timestamp not null default current_timestamp
);
