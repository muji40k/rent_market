
\c rent_market

INSERT INTO users.users (
    id,"token",name,email,"password",modification_date,modification_source
) VALUES
('38df4303-b04c-46f0-b149-2e9ccb4e0dcb'::uuid,'83663176da6ef716d569a4777d6542b4','admin','admin@example.com','admin_psswd',now(),'preset'),
('4e480c88-81dc-4816-a764-99f6075af440'::uuid,'0a74bab2479ab9f18a15119fe5410615','user','user@example.com','user_psswd',now(),'preset'),
('1e54dd94-d67d-46b7-846d-ad80e5f34549'::uuid,'ab0badce0dd2b0585d68d01fedcd3402','renter','renter@example.com','renter_psswd',now(),'preset'),
('bf4b1f09-5282-46bd-ba4f-4739dea0df1c'::uuid,'be83fc2ef182d98973f98e798a12ae9a','storekeeper','sk@example.com','sk_psswd',now(),'preset'),
('75f84237-8aaf-4da0-a9e7-36b634750828'::uuid,'62dc6f9ca19659e87e799f415b78ae91','storekeeper2','sk2@example.com','sk2_psswd',now(),'preset');

INSERT INTO roles.administrators (
    id,user_id,modification_date,modification_source
) VALUES
('b45a6e17-64cc-4cd8-8f56-1293e28087e3'::uuid,'38df4303-b04c-46f0-b149-2e9ccb4e0dcb'::uuid,now(),'preset');

INSERT INTO roles.renters (
    id,user_id,modification_date,modification_source
) VALUES
('04a47ac2-08b7-4139-9c76-2e3fa5bab358'::uuid,'1e54dd94-d67d-46b7-846d-ad80e5f34549'::uuid,now(),'preset');

INSERT INTO roles.storekeepers (
    id,user_id,pick_up_point_id,modification_date,modification_source
) VALUES
('8713f68a-bfd8-4757-a9d9-1f613f6175de'::uuid,'bf4b1f09-5282-46bd-ba4f-4739dea0df1c'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,now(),'preset'),
('6a215830-12bb-4e9f-843a-0d95454c4c5c'::uuid,'75f84237-8aaf-4da0-a9e7-36b634750828'::uuid,'15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,now(),'preset');

