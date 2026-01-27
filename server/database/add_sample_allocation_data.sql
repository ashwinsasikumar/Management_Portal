-- Add sample semesters and course mappings for Course Allocation feature
-- Run this script to populate test data

-- 1. Create semesters for the existing curriculum (ID = 1)
INSERT INTO `normal_cards` (`curriculum_id`, `semester_number`, `visibility`, `card_type`, `status`) VALUES
(1, 1, 'UNIQUE', 'semester', 1),
(1, 2, 'UNIQUE', 'semester', 1),
(1, 3, 'UNIQUE', 'semester', 1),
(1, 4, 'UNIQUE', 'semester', 1),
(1, 5, 'UNIQUE', 'semester', 1),
(1, 6, 'UNIQUE', 'semester', 1),
(1, 7, 'UNIQUE', 'semester', 1),
(1, 8, 'UNIQUE', 'semester', 1);

-- Get the semester IDs that were just created (assuming they start from ID 1)
-- Note: Adjust these IDs if your normal_cards table already has data

-- 2. Link existing courses to semesters
-- Assuming semester IDs are 1-8 based on above inserts
-- Distribute the 5 existing courses across first few semesters

-- Semester 1 courses
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`, `count_towards_limit`) VALUES
(1, 1, 1, 1, 1),  -- CS101 - Programming in C
(1, 1, 5, 1, 1);  -- CS103L - Programming Lab

-- Semester 2 courses  
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`, `count_towards_limit`) VALUES
(1, 2, 2, 1, 1);  -- CS102 - Data Structures

-- Semester 3 courses
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`, `count_towards_limit`) VALUES
(1, 3, 3, 1, 1);  -- IT101 - Web Technologies

-- Semester 4 courses
INSERT INTO `curriculum_courses` (`curriculum_id`, `semester_id`, `course_id`, `status`, `count_towards_limit`) VALUES
(1, 4, 4, 1, 1);  -- EC101 - Digital Electronics

-- Note: The existing course_allocations data should now work
-- because courses are linked to semesters!

-- Verify the data:
-- SELECT * FROM normal_cards WHERE curriculum_id = 1;
-- SELECT * FROM curriculum_courses WHERE curriculum_id = 1;
-- SELECT c.course_code, c.course_name, nc.semester_number 
-- FROM courses c 
-- JOIN curriculum_courses cc ON c.course_id = cc.course_id 
-- JOIN normal_cards nc ON cc.semester_id = nc.id 
-- WHERE cc.curriculum_id = 1;
