-- SELECT *
-- FROM users;

-- SELECT *
-- FROM users
-- WHERE name = 'olga'

-- SELECT *
-- FROM plans;


-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     username VARCHAR,
--     password VARCHAR
-- )

-- CREATE TABLE plans (
--     id SERIAL PRIMARY KEY,
--     title VARCHAR,
--     descriptio VARCHAR,
--     complete BOOLEAN
-- )

-- INSERT INTO users(id, name, password) VALUES (15622, 'Olga', '123');

-- UPDATE users
-- SET name = 'Olga'
-- WHERE id = 11111

-- ALTER TABLE plans
-- ADD COLUMN user_id INTEGER NOT NULL REFERENCES users(id);

-- INSERT INTO plans (title, descriptio, complete, user_id) 
-- VALUES ('My First Plan', 'Description of the plan', false, 15622);

-- SELECT *
-- FROM plans
-- INNER JOIN users ON plans.user_id = users.id