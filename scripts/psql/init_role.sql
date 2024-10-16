\c rent_market

drop schema if exists roles cascade;
create schema roles;

drop table if exists roles.renters;
create table roles.renters
(
    id uuid primary key,
    user_id uuid not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table roles.renters add
    constraint "fkey_renter_user_id"
    foreign key (user_id)
    references users.users(id);

drop table if exists roles.administrators;
create table roles.administrators
(
    id uuid primary key,
    user_id uuid not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table roles.administrators add
    constraint "fkey_administrator_user_id"
    foreign key (user_id)
    references users.users(id);

drop table if exists roles.storekeepers;
create table roles.storekeepers
(
    id uuid primary key,
    user_id uuid not null,
    pick_up_point_id uuid not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table roles.storekeepers add
    constraint "fkey_storekeeper_user_id"
    foreign key (user_id)
    references users.users(id);

alter table roles.storekeepers add
    constraint "fkey_storekeeper_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

