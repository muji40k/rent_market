\c rent_market

drop schema if exists periods cascade;
create schema periods;

drop table if exists periods.periods;
create table periods.periods
(
    id uuid primary key,
    name text not null,
    duration bigint not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

