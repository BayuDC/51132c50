CREATE TABLE "student_assignments"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "student_id" INTEGER NOT NULL,
    "assigment_id" INTEGER NOT NULL,
    "content" JSON NULL,
    "score" DOUBLE PRECISION NULL
);
ALTER TABLE "student_assignments"
ADD CONSTRAINT "student_assignments_student_id_foreign" FOREIGN KEY("student_id") REFERENCES "students"("id") ON DELETE CASCADE;
ALTER TABLE "student_assignments"
ADD CONSTRAINT "student_assignments_assigment_id_foreign" FOREIGN KEY("assigment_id") REFERENCES "assignments"("id") ON DELETE CASCADE;