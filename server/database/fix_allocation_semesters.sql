-- Fix existing course_allocations to have proper semester numbers
-- The existing data has semester values but they need to align with actual semesters

-- Update allocations to use correct semester numbers
-- Current allocations:
-- (1,1,1,'2025-2026',1,'A','Primary') - Course 1, Teacher 1, Sem 1
-- (2,2,2,'2025-2026',3,'A','Primary') - Course 2, Teacher 2, Sem 3
-- (3,3,3,'2025-2026',1,'A','Primary') - Course 3, Teacher 3, Sem 1
-- (4,4,4,'2025-2026',1,'A','Primary') - Course 4, Teacher 4, Sem 1
-- (5,5,1,'2025-2026',1,'A','Primary') - Course 5, Teacher 1, Sem 1
-- (6,5,3,'2025-2026',1,'A','Assistant') - Course 5, Teacher 3, Sem 1

-- These should map to courses correctly based on the curriculum_courses mappings
-- CS101 (course 1) -> Semester 1 -> Teacher 1 (Primary) ✓
-- CS102 (course 2) -> Semester 2 -> Teacher 2 (Primary) - needs update
-- IT101 (course 3) -> Semester 3 -> Teacher 3 (Primary) ✓
-- EC101 (course 4) -> Semester 4 -> Teacher 4 (Primary) - needs update
-- CS103L (course 5) -> Semester 1 -> Teacher 1 (Primary) + Teacher 3 (Assistant) ✓

-- Fix course 2 allocation (should be semester 2, not 3)
UPDATE course_allocations SET semester = 2 WHERE id = 2;

-- Fix course 4 allocation (should be semester 4, not 1)
UPDATE course_allocations SET semester = 4 WHERE id = 4;

-- Verify
-- SELECT ca.id, c.course_code, c.course_name, t.name as teacher, ca.semester, ca.section, ca.role
-- FROM course_allocations ca
-- JOIN courses c ON ca.course_id = c.course_id
-- JOIN teachers t ON ca.teacher_id = t.id
-- WHERE ca.status = 1
-- ORDER BY ca.semester, c.course_code;
