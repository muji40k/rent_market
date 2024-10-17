\c rent_market

drop schema if exists payments cascade;
create schema payments;

drop table if exists payments.methods;
create table payments.methods
(
    id uuid primary key,
    name text not null,
    description text,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

drop table if exists payments.users_methods;
create table payments.users_methods
(
    id uuid primary key,
    pay_method_id uuid not null,
    payer_id text not null,
    user_id uuid not null,
    name text not null,
    priority integer not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table payments.users_methods add
    constraint "fkey_users_methods_method_id"
    foreign key (pay_method_id)
    references payments.methods(id);

alter table payments.users_methods add
    constraint "fkey_users_methods_user_id"
    foreign key (user_id)
    references users.users(id);

drop table if exists payments.payments;
create table payments.payments
(
    id uuid primary key,
    rent_id uuid not null,
    pay_method_id uuid,
    payment_id text,
    period_strat timestamptz not null,
    period_end timestamptz not null,
    currency_id uuid not null,
    value double precision not null,
    status text,
    create_date timestamptz not null,
    payment_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table payments.payments add
    constraint "fkey_payment_rent_id"
    foreign key (rent_id)
    references records.users_rents(id);

alter table payments.payments add
    constraint "fkey_payment_pay_method_id"
    foreign key (pay_method_id)
    references payments.methods(id);


