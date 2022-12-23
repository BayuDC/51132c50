CREATE TABLE "teachers"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "fullname" VARCHAR(255) NOT NULL,
    "user_id" INTEGER NOT NULL
);
ALTER TABLE "teachers"
ADD CONSTRAINT "teachers_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;