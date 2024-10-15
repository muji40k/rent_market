\c rent_market

drop schema if exists provisions cascade;
create schema provisions;

drop table if exists provisions.requests;
create table provisions.requests
(
    id uuid primary key,
    product_id uuid,
    renter_id uuid,
    pick_up_point_id uuid,
    name text,
    description text,
    condition text,
    verification_code text,
    create_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table provisions.requests add
    constraint "fkey_request_product_id"
    foreign key (product_id)
    references products.products(id);

alter table provisions.requests add
    constraint "fkey_request_renter_id"
    foreign key (renter_id)
    references roles.renters(id);

alter table provisions.requests add
    constraint "fkey_request_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

drop table if exists provisions.requests_pay_plans;
create table provisions.requests_pay_plans
(
    id uuid primary key,
    request_id uuid,
    period_id uuid,
    currency_id uuid,
    value double precision,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table provisions.requests_pay_plans add
    constraint "fkey_request_pay_plans_request_id"
    foreign key (request_id)
    references provisions.requests(id);

alter table provisions.requests_pay_plans add
    constraint "fkey_request_pay_plans_period_id"
    foreign key (period_id)
    references periods.periods(id);

alter table provisions.requests_pay_plans add
    constraint "fkey_request_pay_plans_currency_id"
    foreign key (currency_id)
    references currencies.currencies(id);

drop table if exists provisions.revokes;
create table provisions.revokes
(
    id uuid primary key,
    instance_id uuid,
    renter_id uuid,
    pick_up_point_id uuid,
    verification_code text,
    create_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table provisions.revokes add
    constraint "fkey_revoke_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table provisions.revokes add
    constraint "fkey_revoke_renter_id"
    foreign key (renter_id)
    references roles.renters(id);

alter table provisions.revokes add
    constraint "fkey_revoke_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

