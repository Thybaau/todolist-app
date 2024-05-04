-- Create Users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);
--Create Tasks table
CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    content TEXT,
    state BOOLEAN NOT NULL DEFAULT FALSE,
    username VARCHAR(255) REFERENCES users(username) ON DELETE CASCADE
);
