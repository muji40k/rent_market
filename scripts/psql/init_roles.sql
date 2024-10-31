
\c rent_market

revoke all privileges on database rent_market from readonly;

grant usage on schema records to readonly;
grant select on all tables in schema records to readonly;
alter default privileges in schema records grant select on tables to readonly;

grant usage on schema instances to readonly;
grant select on all tables in schema instances to readonly;
alter default privileges in schema instances grant select on tables to readonly;

grant usage on schema products to readonly;
grant select on all tables in schema products to readonly;
alter default privileges in schema products grant select on tables to readonly;

grant usage on schema roles to readonly;
grant select on all tables in schema roles to readonly;
alter default privileges in schema roles grant select on tables to readonly;

grant usage on schema users to readonly;
grant select on all tables in schema users to readonly;
alter default privileges in schema users grant select on tables to readonly;

grant usage on schema pick_up_points to readonly;
grant select on all tables in schema pick_up_points to readonly;
alter default privileges in schema pick_up_points grant select on tables to readonly;

grant usage on schema currencies to readonly;
grant select on all tables in schema currencies to readonly;
alter default privileges in schema currencies grant select on tables to readonly;

grant usage on schema photos to readonly;
grant select on all tables in schema photos to readonly;
alter default privileges in schema photos grant select on tables to readonly;

grant usage on schema addresses to readonly;
grant select on all tables in schema addresses to readonly;
alter default privileges in schema addresses grant select on tables to readonly;

grant usage on schema periods to readonly;
grant select on all tables in schema periods to readonly;
alter default privileges in schema periods grant select on tables to readonly;

grant usage on schema categories to readonly;
grant select on all tables in schema categories to readonly;
alter default privileges in schema categories grant select on tables to readonly;

grant usage on schema payments to readonly;
grant select on all tables in schema payments to readonly;
alter default privileges in schema payments grant select on tables to readonly;

grant usage on schema deliveries to readonly;
grant select on all tables in schema deliveries to readonly;
alter default privileges in schema deliveries grant select on tables to readonly;

grant usage on schema rents to readonly;
grant select on all tables in schema rents to readonly;
alter default privileges in schema rents grant select on tables to readonly;

grant usage on schema provisions to readonly;
grant select on all tables in schema provisions to readonly;
alter default privileges in schema provisions grant select on tables to readonly;

revoke readonly from reader;
grant readonly to reader;

