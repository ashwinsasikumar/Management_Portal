-- Insert academic_details for test students
-- Run this in MySQL Workbench SQL editor

-- First, check if students exist
SELECT COUNT(*) as student_count FROM students WHERE student_id BETWEEN 1001 AND 3999;

-- Delete any existing academic_details for these students first
DELETE FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- Now insert fresh academic_details for all test students
INSERT INTO academic_details (student_id, batch, year, semester, degree_level, section, department, student_category, curriculum_id) VALUES
-- CSE Students (8 students)
(1001, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1002, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1003, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1004, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1005, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),
(1006, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),
(1007, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),
(1008, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),

-- IT Students (5 students)
(2001, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2002, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2003, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2004, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2005, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),

-- ECE Students (4 students)
(3001, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290),
(3002, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290),
(3003, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290),
(3004, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290);

-- Verify the data was inserted
SELECT COUNT(*) as academic_details_count FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- Test the query that the UI uses
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

SELECT 'If you see 8 CSE students above, the data is correct!' as result;
