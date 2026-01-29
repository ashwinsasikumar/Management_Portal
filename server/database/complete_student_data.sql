-- Complete script to insert students and academic details
-- Run this in MySQL Workbench SQL editor

-- Step 1: Clean up existing test data
DELETE FROM student_teacher_mapping WHERE student_id BETWEEN 1001 AND 3999;
DELETE FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;
DELETE FROM students WHERE student_id BETWEEN 1001 AND 3999;

-- Step 2: Insert students first (required for foreign key)
INSERT INTO students (student_id, enrollment_no, student_name, status) VALUES
-- CSE Students (8 students)
(1001, 'CSE2024001', 'Arjun Kumar', 1),
(1002, 'CSE2024002', 'Priya Sharma', 1),
(1003, 'CSE2024003', 'Rahul Verma', 1),
(1004, 'CSE2024004', 'Anjali Patel', 1),
(1005, 'CSE2024005', 'Vikram Singh', 1),
(1006, 'CSE2024006', 'Sneha Reddy', 1),
(1007, 'CSE2024007', 'Karthik Nair', 1),
(1008, 'CSE2024008', 'Divya Menon', 1),

-- IT Students (5 students)
(2001, 'IT2024001', 'Amit Gupta', 1),
(2002, 'IT2024002', 'Neha Joshi', 1),
(2003, 'IT2024003', 'Rajesh Kumar', 1),
(2004, 'IT2024004', 'Pooja Desai', 1),
(2005, 'IT2024005', 'Sanjay Rao', 1),

-- ECE Students (4 students)
(3001, 'ECE2024001', 'Suresh Pillai', 1),
(3002, 'ECE2024002', 'Kavita Iyer', 1),
(3003, 'ECE2024003', 'Manoj Krishnan', 1),
(3004, 'ECE2024004', 'Lakshmi Nair', 1);

-- Verify students inserted
SELECT COUNT(*) as students_inserted FROM students WHERE student_id BETWEEN 1001 AND 3999;

-- Step 3: Now insert academic_details (foreign key will work)
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

-- Verify academic_details inserted
SELECT COUNT(*) as academic_details_inserted FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- Step 4: Test the query that the UI uses for CSE students
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

SELECT 'SUCCESS: If you see 8 CSE students above, refresh the Student-Teacher Mapping page!' as result;
