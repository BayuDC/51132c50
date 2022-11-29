CREATE TABLE IF NOT EXISTS teachers(
    id serial PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);