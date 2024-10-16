\c rent_market

drop schema if exists pick_up_points cascade;
create schema pick_up_points;

drop table if exists pick_up_points.pick_up_points;
create table pick_up_points.pick_up_points
(
    id uuid primary key,
    address_id uuid not null,
    capacity integer not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table pick_up_points.pick_up_points add
    constraint "fkey_pick_up_point_address_id"
    foreign key (address_id)
    references addresses.addresses(id);

drop table if exists pick_up_points.working_hours;
create table pick_up_points.working_hours
(
    id uuid primary key,
    pick_up_point_id uuid not null,
    day integer not null,
    start_time time not null,
    end_time time not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table pick_up_points.working_hours add
    constraint "fkey_working_hours_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

drop table if exists pick_up_points.photos;
create table pick_up_points.photos
(
    id uuid primary key,
    pick_up_point_id uuid not null,
    photo_id uuid not null,
    modification_date timestamptz not null default now(),
    modification_source text not null
);

alter table pick_up_points.photos add
    constraint "fkey_pup_photos_pick_up_point_id"
    foreign key (pick_up_point_id)
    references pick_up_points.pick_up_points(id);

alter table pick_up_points.photos add
    constraint "fkey_pup_photos_photo_id"
    foreign key (photo_id)
    references photos.photos(id);



