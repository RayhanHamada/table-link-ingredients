-- +goose Up

INSERT INTO tm_ingredient (uuid, name, cause_alergy, type, status, created_at)
VALUES
    ('8666d298-517b-45f3-8566-378cb5c8738c', 'Chicken', false, 0, 1, '2025-02-19 09:16:45.151435'),
    ('e97c6a5f-c541-4f3f-84d4-953c3eabe686', 'Pork',    false, 0, 1, '2025-02-19 09:16:45.151435'),
    ('9a3a33cf-7144-4c5d-a0c6-fd8e894a0db5', 'Radish',  false, 2, 1, '2025-02-19 09:16:45.151435'),
    ('b2752259-090e-4e6e-a9a1-2f47d538d833', 'Egg',     true,  1, 1, '2025-02-19 09:16:45.151435')
ON CONFLICT (uuid) DO NOTHING;

INSERT INTO tm_item (uuid, name, price, status, created_at)
VALUES
    ('07419b87-4702-49f9-83aa-f9b489f64b14', 'Chicken Pork',              30000.00, 1, '2025-02-19 09:19:02.37464'),
    ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', 'Chicken Pork with Radish',  35000.00, 1, '2025-02-19 09:19:02.37464'),
    ('7cc760a3-393b-493c-b780-3cfd7afd1cf9', 'Salad Egg',                 20000.00, 1, '2025-02-19 09:19:02.37464')
ON CONFLICT (uuid) DO NOTHING;

INSERT INTO tm_item_ingredient (uuid_item, uuid_ingredient)
VALUES
    ('07419b87-4702-49f9-83aa-f9b489f64b14', '8666d298-517b-45f3-8566-378cb5c8738c'),
    ('07419b87-4702-49f9-83aa-f9b489f64b14', 'e97c6a5f-c541-4f3f-84d4-953c3eabe686'),
    ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', '8666d298-517b-45f3-8566-378cb5c8738c'),
    ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', 'e97c6a5f-c541-4f3f-84d4-953c3eabe686'),
    ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', '9a3a33cf-7144-4c5d-a0c6-fd8e894a0db5'),
    ('7cc760a3-393b-493c-b780-3cfd7afd1cf9', '9a3a33cf-7144-4c5d-a0c6-fd8e894a0db5'),
    ('7cc760a3-393b-493c-b780-3cfd7afd1cf9', 'b2752259-090e-4e6e-a9a1-2f47d538d833')
ON CONFLICT (uuid_item, uuid_ingredient) DO NOTHING;

-- +goose Down

DELETE FROM tm_item_ingredient;
DELETE FROM tm_item;
DELETE FROM tm_ingredient;
