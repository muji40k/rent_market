\c rent_market

drop schema if exists users cascade;
create schema users;

drop table if exists users.users;
create table users.users
(
    id uuid primary key,
    token text,
    name text,
    email text,
    password text,
    modification_date timestamptz not null default now(),
    modification_source text
);

-- alter table categories.categories add
--     constraint "fkey_category_parent_id"
--     foreign key (parent_id)
--     references categories.categories(id);

drop table if exists users.favorite_pick_up_points;
create table users.favorite_pick_up_points
(
    id uuid primary key,
    user_id uuid,
    pick_up_point_id uuid,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table users.favorite_pick_up_points add
    constraint "fkey_favorite_pick_up_point_user_id"
    foreign key (user_id)
    references users.users(id);

alter table users.favorite_pick_up_points add
    constraint "fkey_favorite_pick_up_point_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

drop table if exists users.profiles;
create table users.profiles
(
    id uuid primary key,
    user_id uuid,
    name text,
    surname text,
    patronymic text,
    birth_date timestamptz,
    photo_id uuid,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table users.profiles add
    constraint "fkey_profile_user_id"
    foreign key (user_id)
    references users.users(id);

alter table users.profiles add
    constraint "fkey_profile_photo_id"
    foreign key (photo_id)
    references photos.photos(id);


