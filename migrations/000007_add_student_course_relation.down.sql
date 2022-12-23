ALTER TABLE "student_courses" DROP CONSTRAINT "student_courses_student_id_foreign";
ALTER TABLE "student_courses" DROP CONSTRAINT "student_courses_course_id_foreign";
DROP TABLE "student_courses";