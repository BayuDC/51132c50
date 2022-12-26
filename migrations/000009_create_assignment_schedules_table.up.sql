CREATE TABLE "assignment_schedules"(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "assignment_id" INTEGER NOT NULL,
    "group_id" INTEGER NOT NULL,
    "open_at" TIMESTAMP WITHOUT TIME ZONE NULL,
    "close_at" TIMESTAMP WITHOUT TIME ZONE NULL
);
ALTER TABLE "assignment_schedules"
ADD CONSTRAINT "assignment_schedules_assignment_id_foreign" FOREIGN KEY("assignment_id") REFERENCES "assignments"("id") ON DELETE CASCADE;
ALTER TABLE "assignment_schedules"
ADD CONSTRAINT "assignment_schedules_group_id_foreign" FOREIGN KEY("group_id") REFERENCES "groups"("id") ON DELETE CASCADE;