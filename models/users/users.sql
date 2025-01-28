CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);
