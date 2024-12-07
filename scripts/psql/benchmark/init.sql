
\i /scripts/users.sql
\i /scripts/roles.sql
\i /scripts/init.sql

\c rent_market
\i /testdata/categories.sql
\copy products.products(id, name, category_id, description, modification_source) from '/testdata/bench_products.csv' delimiter ',' csv
\copy products.characteristics(id, product_id, name, value, modification_source) from '/testdata/bench_product_chars.csv' delimiter ',' csv

