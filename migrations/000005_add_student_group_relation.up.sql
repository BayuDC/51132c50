ALTER TABLE "students"
ADD CONSTRAINT "students_group_id_foreign" FOREIGN KEY("group_id") REFERENCES "groups"("id") ON DELETE RESTRICT;