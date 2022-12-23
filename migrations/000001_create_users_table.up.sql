CREATE TABLE "users"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "username" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NULL,
    "role" VARCHAR(255) CHECK ("role" IN('admin', 'student', 'teacher')) NOT NULL
);