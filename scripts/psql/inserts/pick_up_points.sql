
\c rent_market

INSERT INTO addresses.addresses (
    id,country,city,street,house,flat,modification_date,modification_source
) VALUES
('a05526b9-3f58-447b-bc79-7751e549ade0'::uuid,'Россия','Москва','Примерная','13',NULL,now(),'preset'),
('682a48b2-d137-4a40-97a5-cf6ba3613e9d'::uuid,'Россия','Москва','Домашняя','1',NULL,now(),'preset'),
('b28bb97a-3c82-4b35-a3c4-964b3573d681'::uuid,'Россия','Калининград','Складская','43',NULL,now(),'preset');

INSERT INTO pick_up_points.pick_up_points (
    id,address_id,capacity,modification_date,modification_source
) VALUES
('bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,'b28bb97a-3c82-4b35-a3c4-964b3573d681'::uuid,1337,now(),'preset'),
('15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,'682a48b2-d137-4a40-97a5-cf6ba3613e9d'::uuid,1235,now(),'preset');

INSERT INTO pick_up_points.working_hours (
    id,pick_up_point_id,"day",start_time,end_time,modification_date,modification_source
) VALUES
('1790f6dd-0ff5-42f8-a407-06729af319f7'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,1,'08:00:00','21:00:00',now(),'preset'),
('71f6964c-39f8-4c23-8003-0f5e51d0531e'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,2,'08:00:00','21:00:00',now(),'preset'),
('e7d4986f-ee00-44ce-a77c-021a3e3e9ce2'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,3,'08:00:00','21:00:00',now(),'preset'),
('50d3f1c2-bc93-4822-8090-9b94c5c1c51e'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,4,'08:00:00','21:00:00',now(),'preset'),
('e44b846d-7544-420a-969f-c6d300428b67'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,5,'08:00:00','21:00:00',now(),'preset'),
('2a6e5978-f2d7-43f0-b8e0-5fa0e2faabb6'::uuid,'bba7e35c-93e7-4f95-9781-0b8f364d7f85'::uuid,6,'10:00:00','17:00:00',now(),'preset'),
('c4b0f41a-df5b-4414-bdef-c7d967e75889'::uuid,'15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,1,'08:00:00','20:00:00',now(),'preset'),
('b5c122ec-a79e-4820-824d-60693894d9a6'::uuid,'15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,2,'08:00:00','20:00:00',now(),'preset'),
('57922f15-e6d6-4e70-899b-a1c2633e391c'::uuid,'15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,3,'08:00:00','20:00:00',now(),'preset'),
('343bb18c-1148-493e-a791-46c7333ca21c'::uuid,'15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,4,'08:00:00','20:00:00',now(),'preset'),
('83b40d87-d978-4e4c-a3ac-4308fc0a5722'::uuid,'15a90e2c-4e6f-479b-bbf4-1b3f2c2b99b2'::uuid,5,'08:00:00','18:00:00',now(),'preset');



