CREATE TABLE IF NOT EXISTS products (
    id SERIAL   PRIMARY KEY,
    title       TEXT,
    description TEXT,
    price       DECIMAL,
    category_id INTEGER,
    CONSTRAINT  fk_category
    FOREIGN KEY(category_id)
        REFERENCES categories(id)
        ON DELETE SET NULL
);