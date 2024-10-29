
\c rent_market

INSERT INTO periods.periods (
    id,"name",duration,modification_date,modification_source
) VALUES
('e6183f61-221d-4eaf-8cad-298edf681347'::uuid,'day',86400000000000,now(),'preset'),
('4a238a70-05f5-463b-a1ed-965e0cc05519'::uuid,'week',604800000000000,now(),'preset'),
('47494bab-a13f-4f6c-aea5-2080ee3b7b1d'::uuid,'month',2592000000000000,now(),'preset'),
('5d8ef1b5-f5d7-4dfb-af2d-df074283da6e'::uuid,'quarter',7776000000000000,now(),'preset'),
('a78fd7e4-7d22-4551-ac96-3aa9f0793cf6'::uuid,'half',15552000000000000,now(),'preset'),
('64ead809-ccc2-4dda-be6c-529b76515e9a'::uuid,'year',31104000000000000,now(),'preset');

