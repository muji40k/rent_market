\c rent_market

drop schema if exists photos cascade;
create schema photos;

drop table if exists photos.photos;
create table photos.photos
(
    id uuid primary key,
    placeholder text not null,
    description text not null,
    path text not null,
    mime text not null,
    date timestamptz not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

drop table if exists photos.temp;
create table photos.temp
(
    id uuid primary key,
    placeholder text not null,
    description text not null,
    path text,
    mime text not null,
    date timestamptz not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

