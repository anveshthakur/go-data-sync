\connect sourcedb;

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    quantity INTEGER,
    price DECIMAL(10, 2)
);

INSERT INTO items (name, description, quantity, price) VALUES
('Item 1', 'Description for Item 1', 10, 100.00),
('Item 2', 'Description for Item 2', 20, 200.00),
('Item 3', 'Description for Item 3', 30, 300.00),
('Item 4', 'Description for Item 4', 40, 400.00),
('Item 5', 'Description for Item 5', 50, 500.00);

DO $$
BEGIN
    FOR i IN 6..200 LOOP
        EXECUTE format(
            'INSERT INTO items (name, description, quantity, price) VALUES (%L, %L, %s, %s);',
            'Item ' || i, 'Description for Item ' || i, i * 10, i * 100.00
        );
    END LOOP;
END $$;

\connect targetdb;

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    quantity INTEGER,
    price DECIMAL(10, 2)
);

INSERT INTO items (name, description, quantity, price) VALUES
('Item 1', 'Description for Item 1', 10, 100.00),
('Item 2', 'Description for Item 2', 20, 200.00),
('Item 3', 'Description for Item 3', 30, 300.00),
('Item 4', 'Description for Item 4', 40, 400.00),
('Item 5', 'Description for Item 5', 50, 500.00);

DO $$
BEGIN
    FOR i IN 6..200 LOOP
        EXECUTE format(
            'INSERT INTO items (name, description, quantity, price) VALUES (%L, %L, %s, %s);',
            'Item ' || i, 'Description for Item ' || i, i * 10, i * 100.00 + 1
        );
    END LOOP;
END $$;