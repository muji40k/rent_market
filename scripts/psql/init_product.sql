\c rent_market

drop schema if exists products cascade;
create schema products;

drop table if exists products.products;
create table products.products
(
    id uuid primary key,
    name text,
    category_id uuid,
    description text,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table products.products add
    constraint "fkey_product_category_id"
    foreign key (category_id)
    references categories.categories(id);

drop table if exists products.characteristics;
create table products.characteristics
(
    id uuid primary key,
    product_id uuid,
    name text,
    value text,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table products.characteristics add
    constraint "fkey_characteristics_product_id"
    foreign key (product_id)
    references products.products(id);

drop table if exists products.photos;
create table products.photos
(
    id uuid primary key,
    product_id uuid,
    photo_id uuid,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table products.photos add
    constraint "fkey_product_photos_product_id"
    foreign key (product_id)
    references products.products(id);

alter table products.photos add
    constraint "fkey_product_photos_photo_id"
    foreign key (photo_id)
    references photos.photos(id);

