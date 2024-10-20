\c rent_market

drop schema if exists products cascade;
create schema products;

drop function if exists to_tsvector_mutilang;
create function to_tsvector_multilang(text)
returns tsvector
as $$
    select to_tsvector('english', $1) || to_tsvector('russian', $1)
$$ language sql immutable;

drop table if exists products.products;
create table products.products
(
    id uuid primary key,
    name text not null,
    category_id uuid not null,
    description text not null,
    modification_date timestamptz not null default now(),
    modification_source text not null,
    ts tsvector generated always as (to_tsvector_multilang(name || ' ' || description)) stored
);

alter table products.products add
    constraint "fkey_product_category_id"
    foreign key (category_id)
    references categories.categories(id);

drop table if exists products.characteristics;
create table products.characteristics
(
    id uuid primary key,
    product_id uuid not null,
    name text not null,
    value text not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table products.characteristics add
    constraint "fkey_characteristics_product_id"
    foreign key (product_id)
    references products.products(id);

drop table if exists products.photos;
create table products.photos
(
    id uuid primary key,
    product_id uuid not null,
    photo_id uuid not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table products.photos add
    constraint "fkey_product_photos_product_id"
    foreign key (product_id)
    references products.products(id);

alter table products.photos add
    constraint "fkey_product_photos_photo_id"
    foreign key (photo_id)
    references photos.photos(id);

