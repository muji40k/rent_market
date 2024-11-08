
\c rent_market

INSERT INTO payments.methods (
    id, "name", description, modification_date, modification_source
) VALUES (
    'a2e7d44f-0af0-4b9b-bd70-65d3397d5ad9'::uuid, 'dummy',
    'Еще одна заглушка для возвращения истины', now(), 'psql'
);

