\c rent_market

drop schema if exists addresses cascade;
create schema addresses;

drop table if exists addresses.addresses;
create table addresses.addresses
(
    id uuid primary key,
    country text not null,
    city text not null,
    street text not null,
    house text not null,
    flat text,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

