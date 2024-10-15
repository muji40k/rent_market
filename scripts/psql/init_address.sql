\c rent_market

drop schema if exists addresses cascade;
create schema addresses;

drop table if exists addresses.addresses;
create table addresses.addresses
(
    id uuid primary key,
    country text,
    city text,
    street text,
    house text,
    flat text,
    modification_date timestamptz not null default now(),
    modification_source text
);

