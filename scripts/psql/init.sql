\c postgres

drop database if exists rent_market;
create database rent_market;

\i /scripts/init_category.sql
\i /scripts/init_period.sql
\i /scripts/init_address.sql
\i /scripts/init_photo.sql
\i /scripts/init_currency.sql

\i /scripts/init_pick_up_point.sql

\i /scripts/init_user.sql
\i /scripts/init_role.sql

\i /scripts/init_product.sql
\i /scripts/init_instance.sql

\i /scripts/init_records.sql

\i /scripts/init_provision_requests.sql
\i /scripts/init_rent_requests.sql
\i /scripts/init_delivery.sql

\i /scripts/init_payment.sql

