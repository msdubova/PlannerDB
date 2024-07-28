CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL
);

CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    descriptio TEXT NOT NULL,
    complete BOOLEAN NOT NULL
);