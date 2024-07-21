

-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR
-- );


-- SELECT *
-- FROM users WHERE name='Lena';

-- INSERT INTO users(name) VALUES ('John')

-- UPDATE users
-- SET name = 'Gena'
-- WHERE id = 3

-- SELECT *
-- FROM users

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id)
)