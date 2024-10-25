\c postgres

drop database if exists authentication;
create database authentication;

\c authentication

drop table if exists public.sessions;
create table public.sessions (
    access_token text primary key,
    access_valid_to timestamptz not null,
    renew_token text not null,
    renew_valid_to timestamptz not null,
    token text not null
);

