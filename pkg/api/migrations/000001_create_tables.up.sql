CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          name TEXT
);

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username TEXT,
                                     email TEXT,
                                     password_hash TEXT
);

CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        title TEXT,
                                        description TEXT,
                                        price DECIMAL,
                                        category_id INTEGER,
                                        CONSTRAINT fk_category
                                        FOREIGN KEY(category_id)
    REFERENCES categories(id)
    ON DELETE SET NULL
    );

CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id INTEGER,
                                      total_price DECIMAL,
                                      status TEXT,
                                      CONSTRAINT fk_user
                                      FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    );
