\c rent_market

drop schema if exists deliveries cascade;
create schema deliveries;

drop table if exists deliveries.companies;
create table deliveries.companies
(
    id uuid primary key,
    name text,
    site text,
    phone_bumber text,
    description text,
    modification_date timestamptz not null default now(),
    modification_source text
);

drop table if exists deliveries.deliveries;
create table deliveries.deliveries
(
    id uuid primary key,
    company_id uuid,
    instance_id uuid,
    from_id uuid,
    to_id uuid,
    delivery_id text,
    scheduled_begin_date timestamptz,
    actual_begin_date timestamptz,
    scheduled_end_date timestamptz,
    actual_end_date timestamptz,
    verification_code text,
    create_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table deliveries.deliveries add
    constraint "fkey_delivery_company_id"
    foreign key (company_id)
    references deliveries.companies(id);

alter table deliveries.deliveries add
    constraint "fkey_delivery_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table deliveries.deliveries add
    constraint "fkey_delivery_from_id"
    foreign key (from_id)
    references pick_up_points.pick_up_points(id);

alter table deliveries.deliveries add
    constraint "fkey_delivery_to_id"
    foreign key (to_id)
    references pick_up_points.pick_up_points(id);

