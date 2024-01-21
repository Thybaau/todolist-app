--Create Tasks table
CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    content TEXT,
    state BOOLEAN NOT NULL DEFAULT FALSE
);