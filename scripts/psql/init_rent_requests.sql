\c rent_market

drop schema if exists rents cascade;
create schema rents;

drop table if exists rents.requests;
create table rents.requests
(
    id uuid primary key,
    instance_id uuid not null,
    user_id uuid not null,
    pick_up_point_id uuid not null,
    payment_period_id uuid not null,
    verification_code text not null,
    create_date timestamptz not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table rents.requests add
    constraint "fkey_rent_requests_isntance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table rents.requests add
    constraint "fkey_rent_request_user_id"
    foreign key (user_id)
    references users.users(id);

alter table rents.requests add
    constraint "fkey_rent_request_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

alter table rents.requests add
    constraint "fkey_rent_period_id"
    foreign key (payment_period_id)
    references periods.periods(id);

drop table if exists rents.returns;
create table rents.returns
(
    id uuid primary key,
    instance_id uuid not null,
    user_id uuid not null,
    pick_up_point_id uuid not null,
    rent_end_date timestamptz not null,
    verification_code text not null,
    create_date timestamptz not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table rents.returns add
    constraint "fkey_rent_returns_isntance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table rents.returns add
    constraint "fkey_rent_returns_user_id"
    foreign key (user_id)
    references users.users(id);

alter table rents.returns add
    constraint "fkey_rent_returns_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

