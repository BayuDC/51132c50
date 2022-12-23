ALTER TABLE "assignment_schedules" DROP CONSTRAINT "assignment_schedules_assignment_id_foreign";
ALTER TABLE "assignment_schedules" DROP CONSTRAINT "assignment_schedules_group_id_foreign";
DROP TABLE "assignment_schedules";