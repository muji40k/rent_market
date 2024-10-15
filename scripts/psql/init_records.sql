\c rent_market

drop schema if exists records cascade;
create schema records;

drop table if exists records.users_rents;
create table records.users_rents
(
    id uuid primary key,
    user_id uuid,
    instance_id uuid,
    start_date timestamptz,
    end_date timestamptz,
    payment_period_id uuid,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table records.users_rents add
    constraint "fkey_user_rents_user_id"
    foreign key (user_id)
    references users.users(id);

alter table records.users_rents add
    constraint "fkey_user_rents_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table records.users_rents add
    constraint "fkey_user_rents_period_id"
    foreign key (payment_period_id)
    references periods.periods(id);

drop table if exists records.renters_instances;
create table records.renters_instances
(
    id uuid primary key,
    renter_id uuid,
    instance_id uuid,
    start_date timestamptz,
    end_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table records.renters_instances add
    constraint "fkey_renter_instances_renter_id"
    foreign key (renter_id)
    references roles.renters(id);

alter table records.renters_instances add
    constraint "fkey_renter_instances_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

drop table if exists records.pick_up_points_instances;
create table records.pick_up_points_instances
(
    id uuid primary key,
    pick_up_point_id uuid,
    instance_id uuid,
    in_date timestamptz,
    out_date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table records.pick_up_points_instances add
    constraint "fkey_pick_up_point_instances_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

alter table records.pick_up_points_instances add
    constraint "fkey_pick_up_point_instances_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

