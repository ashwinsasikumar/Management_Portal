-- Fix academic_details with existing curriculum_id
-- Run this in MySQL Workbench SQL editor

-- Step 1: Check what curriculum IDs exist
SELECT id, name, academic_year FROM curriculum ORDER BY id;

-- Step 2: Use the first curriculum ID (or NULL if curriculum_id is optional)
-- Delete existing academic_details for test students
DELETE FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- Step 3: Insert academic_details WITHOUT curriculum_id (set it to NULL)
INSERT INTO academic_details (student_id, batch, year, semester, degree_level, section, department, student_category) VALUES
-- CSE Students (8 students)
(1001, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular'),
(1002, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular'),
(1003, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular'),
(1004, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular'),
(1005, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular'),
(1006, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular'),
(1007, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular'),
(1008, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular'),

-- IT Students (5 students)
(2001, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular'),
(2002, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular'),
(2003, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular'),
(2004, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular'),
(2005, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular'),

-- ECE Students (4 students)
(3001, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular'),
(3002, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular'),
(3003, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular'),
(3004, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular');

-- Verify the data was inserted
SELECT COUNT(*) as academic_details_count FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- Test the query that the UI uses for CSE students
SELECT 
    s.student_id,
    s.enrollment_no,
    s.student_name,
    ad.department,
    ad.year
FROM students s
INNER JOIN academic_details ad ON s.student_id = ad.student_id
INNER JOIN departments d ON ad.department = d.department_name
WHERE d.id = 1 AND ad.year = 1 AND s.status = 1
ORDER BY s.enrollment_no;

SELECT 'SUCCESS: If you see 8 CSE students above, refresh the UI!' as result;
