CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    token_hash TEXT UNIQUE NOT NULL 
);