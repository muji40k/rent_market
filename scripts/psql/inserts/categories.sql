
\c rent_market

INSERT INTO categories.categories (
    id, "name", parent_id, modification_date, modification_source
) VALUES
('2c1a5afd-043f-42c3-8563-1ff24f0adea2'::uuid, 'Одежда', null, now(), 'preset'),
('e19dfa5d-902a-484c-95fb-4bf0818d4420'::uuid, 'Готовый костюм', '2c1a5afd-043f-42c3-8563-1ff24f0adea2'::uuid, now(), 'preset'),
('80c172b1-ae1e-4baa-9cd3-a47dd2047ca5'::uuid, 'Верхняя', '2c1a5afd-043f-42c3-8563-1ff24f0adea2'::uuid, now(), 'preset'),
('e08c9037-8fc9-424a-837e-eacfbf1c637d'::uuid, 'Брюки', '2c1a5afd-043f-42c3-8563-1ff24f0adea2'::uuid, now(), 'preset'),
('bebe188c-da95-4dfd-a431-bed64a204e7f'::uuid, 'Обувь', '2c1a5afd-043f-42c3-8563-1ff24f0adea2'::uuid, now(), 'preset'),
('d9114bdc-e10c-4b60-8939-5d10c833daa8'::uuid, 'Электроника', null, now(), 'preset'),
('923297f7-3dc3-4d02-a4ef-f82dd4ce1fb9'::uuid, 'Компьютеры', 'd9114bdc-e10c-4b60-8939-5d10c833daa8'::uuid, now(), 'preset'),
('a12b94b6-9d01-406f-868c-fc2e8efb1524'::uuid, 'Переферия', '923297f7-3dc3-4d02-a4ef-f82dd4ce1fb9'::uuid, now(), 'preset'),
('03d35033-c447-41fe-90a1-bdfc80536991'::uuid, 'Клавиатуры', 'a12b94b6-9d01-406f-868c-fc2e8efb1524'::uuid, now(), 'preset'),
('3973e21f-a423-4720-bd91-ec5ff9b90bd8'::uuid, 'Комплектующие', '923297f7-3dc3-4d02-a4ef-f82dd4ce1fb9'::uuid, now(), 'preset'),
('4fde0ac7-0d81-421e-97c9-d6d0957c2f7d'::uuid, 'Видеокарты', '3973e21f-a423-4720-bd91-ec5ff9b90bd8'::uuid, now(), 'preset'),
('e373fc89-1023-4237-bb4f-4f8f3b653b00'::uuid, 'Фото/Видео техника', 'd9114bdc-e10c-4b60-8939-5d10c833daa8'::uuid, now(), 'preset'),
('ff406572-79c2-4040-b2a8-e0605ab9cf57'::uuid, 'Фотоаппараты', 'e373fc89-1023-4237-bb4f-4f8f3b653b00'::uuid, now(), 'preset'),
('32e76612-1cbd-443a-a54a-7be810f87c1d'::uuid, 'Видеокамеры', 'e373fc89-1023-4237-bb4f-4f8f3b653b00'::uuid, now(), 'preset'),
('9b8e6e20-0d57-4415-a2df-9504e85492ac'::uuid, 'Объективы', 'e373fc89-1023-4237-bb4f-4f8f3b653b00'::uuid, now(), 'preset');

