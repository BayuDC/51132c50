CREATE TABLE "students"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "fullname" VARCHAR(255) NOT NULL,
    "user_id" INTEGER NOT NULL,
    "group_id" INTEGER NULL
);
ALTER TABLE "students"
ADD CONSTRAINT "students_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;