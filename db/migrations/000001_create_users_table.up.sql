DROP TYPE IF EXISTS role;
CREATE TYPE role as ENUM('admin', 'teacher', 'student');
CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NULL,
    role role NOT NULL
);