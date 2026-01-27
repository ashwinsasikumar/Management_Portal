-- Complete setup script for Course Allocation feature
-- This adds all necessary test data to make the feature work

-- ========================================
-- 1. CREATE SEMESTERS
-- ========================================
-- Add 8 semesters for the existing curriculum (ID = 1: R2025-2026 B.E CSE)
INSERT INTO `normal_cards` (`curriculum_id`, `semester_number`, `visibility`, `card_type`, `status`) VALUES
(1, 1, 'UNIQUE', 'semester', 1),
(1, 2, 'UNIQUE', 'semester', 1),
(1, 3, 'UNIQUE', 'semester', 1),
(1, 4, 'UNIQUE', 'semester', 1),
(1, 5, 'UNIQUE', 'semester', 1),
(1, 6, 'UNIQUE', 'semester', 1),
(1, 7, 'UNIQUE', 'semester', 1),
(1, 8, 'UNIQUE', 'semester', 1);

-- ========================================
-- 2. LINK COURSES TO SEMESTERS
-- ========================================
-- Map existing courses to semesters (curriculum_courses table)
-- Existing courses:
-- 1 = CS101 - Programming in C
-- 2 = CS102 - Data Structures  
-- 3 = IT101 - Web Technologies
-- 4 = EC101 - Digital Electronics
-- 5 = CS103L - Programming Lab

-- Get semester IDs (they should be 1-8 if normal_cards was empty)
-- Semester 1
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`) VALUES
(1, 1, 1, 1),  -- CS101
(1, 1, 5, 1);  -- CS103L (Lab)

-- Semester 2
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`) VALUES
(1, 2, 2, 1);  -- CS102

-- Semester 3
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`) VALUES
(1, 3, 3, 1);  -- IT101

-- Semester 4
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`) VALUES
(1, 4, 4, 1);  -- EC101

-- ========================================
-- 3. FIX EXISTING ALLOCATIONS
-- ========================================
-- Update existing course_allocations to match correct semesters
UPDATE course_allocations SET semester = 2 WHERE id = 2;  -- CS102 is in semester 2
UPDATE course_allocations SET semester = 4 WHERE id = 4;  -- EC101 is in semester 4

-- ========================================
-- VERIFICATION QUERIES
-- ========================================
-- Run these to verify everything is set up correctly:

-- Check semesters
-- SELECT id, curriculum_id, semester_number, card_type FROM normal_cards WHERE curriculum_id = 1 ORDER BY semester_number;

-- Check course-semester mappings
-- SELECT cc.id, nc.semester_number, c.course_code, c.course_name 
-- FROM curriculum_courses cc
-- JOIN normal_cards nc ON cc.semester_id = nc.id
-- JOIN courses c ON cc.course_id = c.course_id
-- WHERE cc.curriculum_id = 1
-- ORDER BY nc.semester_number;

-- Check allocations
-- SELECT ca.id, nc.semester_number, c.course_code, t.name as teacher, ca.section, ca.role
-- FROM course_allocations ca
-- JOIN courses c ON ca.course_id = c.course_id
-- JOIN teachers t ON ca.teacher_id = t.id
-- LEFT JOIN curriculum_courses cc ON ca.course_id = cc.course_id
-- LEFT JOIN normal_cards nc ON cc.semester_id = nc.id
-- WHERE ca.status = 1
-- ORDER BY nc.semester_number, c.course_code;

-- ========================================
-- EXPECTED RESULTS
-- ========================================
-- After running this script, you should have:
-- - 8 semesters created
-- - 5 courses distributed across 4 semesters
-- - 6 teacher allocations properly linked
-- 
-- Now the Course Allocation page will show:
-- - Curriculum dropdown with "R2025-2026 B.E CSE"
-- - Semester dropdown with "Semester 1" through "Semester 8"
-- - When selecting Semester 1, you'll see CS101 and CS103L
-- - When selecting Semester 2, you'll see CS102
-- - etc.
