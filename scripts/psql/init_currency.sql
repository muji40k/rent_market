\c rent_market

drop schema if exists currencies cascade;
create schema currencies;

drop table if exists currencies.currencies;
create table currencies.currencies
(
    id uuid primary key,
    name text,
    modification_date timestamptz not null default now(),
    modification_source text
);

