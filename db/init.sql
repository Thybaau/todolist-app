CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(100),
    content VARCHAR(600),
    statu BOOLEAN NOT NULL
);