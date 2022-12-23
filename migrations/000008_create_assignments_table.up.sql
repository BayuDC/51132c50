CREATE TABLE "assignments"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "course_id" INTEGER NOT NULL,
    "type" VARCHAR(255) CHECK ("type" IN('empty', 'files')) NOT NULL
);
ALTER TABLE "assignments"
ADD CONSTRAINT "assignments_course_id_foreign" FOREIGN KEY("course_id") REFERENCES "courses"("id") ON DELETE CASCADE;