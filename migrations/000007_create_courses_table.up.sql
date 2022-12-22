CREATE TABLE IF NOT EXISTS courses(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    teacher_id INTEGER NULL,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE
    SET NULL
)