CREATE TABLE "student_courses"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "student_id" INTEGER NOT NULL,
    "course_id" INTEGER NOT NULL
);
ALTER TABLE "student_courses"
ADD CONSTRAINT "student_courses_student_id_foreign" FOREIGN KEY("student_id") REFERENCES "students"("id") ON DELETE CASCADE;
ALTER TABLE "student_courses"
ADD CONSTRAINT "student_courses_course_id_foreign" FOREIGN KEY("course_id") REFERENCES "courses"("id") ON DELETE CASCADE;