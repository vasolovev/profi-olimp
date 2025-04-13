-- Drop indexes
DROP INDEX IF EXISTS idx_students_group_id;
DROP INDEX IF EXISTS idx_groups_parent_id;
DROP INDEX IF EXISTS idx_students_name;
DROP INDEX IF EXISTS idx_groups_name;

-- Drop tables
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS groups;