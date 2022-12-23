CREATE TABLE "courses"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "teacher_id" INTEGER NULL
);
ALTER TABLE "courses"
ADD CONSTRAINT "courses_teacher_id_foreign" FOREIGN KEY("teacher_id") REFERENCES "teachers"("id") ON DELETE
SET NULL;