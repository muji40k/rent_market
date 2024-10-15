\c rent_market

drop schema if exists payments cascade;
create schema payments;

drop table if exists payments.methods;
create table payments.methods
(
    id uuid primary key,
    name text,
    description text,
    modification_date timestamptz not null default now(),
    modification_source text
);

drop table if exists payments.users_methods;
create table payments.users_methods
(
    id uuid primary key,
    pay_method_id uuid,
    payer_id text,
    user_id uuid,
    name text,
    priority integer,
    modification_date timestamptz not null default now(),
    modification_source text
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
    rent_id uuid,
    pay_method_id uuid,
    payment_id text,
    period_strat timestamptz,
    period_end timestamptz,
    currency_id uuid,
    value double precision,
    status text,
    create_date timestamptz,
    payment_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table payments.payments add
    constraint "fkey_payment_rent_id"
    foreign key (rent_id)
    references records.users_rents(id);

alter table payments.payments add
    constraint "fkey_payment_pay_method_id"
    foreign key (pay_method_id)
    references payments.methods(id);


