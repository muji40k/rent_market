
\c rent_market

INSERT INTO deliveries.companies (
    id, "name", site, phone_bumber, description, modification_date,
    modification_source
) VALUES (
    '44072c0a-e312-452a-aa20-69429d7b950d'::uuid, 'dummy', 'localhost:42069',
    '8-800-555-35-35',
    'Заглушка, которая постоянно возвращает шаблонный положительный результат',
    now(), 'psql'
);

