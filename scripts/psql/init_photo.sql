\c rent_market

drop schema if exists photos cascade;
create schema photos;

drop table if exists photos.photos;
create table photos.photos
(
    id uuid primary key,
    placeholder text,
    description text,
    path text,
    mime text,
    date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

drop table if exists photos.temp;
create table photos.temp
(
    id uuid primary key,
    placeholder text,
    description text,
    path text,
    mime text,
    date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

