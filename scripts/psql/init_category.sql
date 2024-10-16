\c rent_market

drop schema if exists categories cascade;
create schema categories;

drop table if exists categories.categories;
create table categories.categories
(
    id uuid primary key,
    name text not null,
    parent_id uuid,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table categories.categories add
    constraint "fkey_category_parent_id"
    foreign key (parent_id)
    references categories.categories(id);

