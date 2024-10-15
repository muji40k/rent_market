\c rent_market

drop schema if exists instances cascade;
create schema instances;

drop table if exists instances.instances;
create table instances.instances
(
    id uuid primary key,
    product_id uuid,
    name text,
    description text,
    condition text,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table instances.instances add
    constraint "fkey_instance_product_id"
    foreign key (product_id)
    references products.products(id);

drop table if exists instances.pay_plans;
create table instances.pay_plans
(
    id uuid primary key,
    instance_id uuid,
    period_id uuid,
    currency_id uuid,
    price double precision,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table instances.pay_plans add
    constraint "fkey_pay_plans_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table instances.pay_plans add
    constraint "fkey_pay_plans_period_id"
    foreign key (period_id)
    references periods.periods(id);

alter table instances.pay_plans add
    constraint "fkey_pay_plans_currency_id"
    foreign key (currency_id)
    references currencies.currencies(id);

drop table if exists instances.photos;
create table instances.photos
(
    id uuid primary key,
    instance_id uuid,
    photo_id uuid,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table instances.photos add
    constraint "fkey_instance_photos_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table instances.photos add
    constraint "fkey_instance_photos_photo_id"
    foreign key (photo_id)
    references photos.photos(id);

drop table if exists instances.reviews;
create table instances.reviews
(
    id uuid primary key,
    instance_id uuid,
    user_id uuid,
    content text,
    rating double precision,
    date timestamptz,
    modification_date timestamptz not null default now(),
    modification_source text
);

alter table instances.reviews add
    constraint "fkey_review_instance_id"
    foreign key (instance_id)
    references instances.instances(id);

alter table instances.reviews add
    constraint "fkey_review_user_id"
    foreign key (user_id)
    references users.users(id);

