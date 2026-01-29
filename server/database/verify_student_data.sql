-- Verify student data for debugging
-- Run these queries one by one to check the data

-- 1. Check if students exist
SELECT COUNT(*) as total_students FROM students WHERE student_id BETWEEN 1001 AND 3999;
SELECT * FROM students WHERE student_id BETWEEN 1001 AND 3999 LIMIT 5;

-- 2. Check if academic_details exist
SELECT COUNT(*) as total_academic_details FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;
SELECT * FROM academic_details WHERE student_id BETWEEN 1001 AND 3999 LIMIT 5;

-- 3. Check departments table
SELECT * FROM departments;

-- 4. Check if department names match exactly
SELECT DISTINCT ad.department FROM academic_details ad WHERE student_id BETWEEN 1001 AND 3999;
SELECT department_name FROM departments WHERE id IN (1, 2, 3);

-- 5. Test the EXACT query used by the API for CSE students (department_id=1, year=1)
SELECT 
    s.student_id,
    s.enrollment_no,
    s.student_name,
    ad.department,
    ad.year,
    ad.section,
    d.id as dept_id,
    d.department_name
FROM students s
INNER JOIN academic_details ad ON s.student_id = ad.student_id
INNER JOIN departments d ON ad.department = d.department_name
WHERE d.id = 1 AND ad.year = 1 AND s.status = 1
ORDER BY s.enrollment_no;

-- 6. Check if academic_year filter matters (remove it to test)
SELECT 
    s.student_id,
    s.enrollment_no,
    s.student_name,
    ad.department,
    ad.year
FROM students s
INNER JOIN academic_details ad ON s.student_id = ad.student_id
INNER JOIN departments d ON ad.department = d.department_name
WHERE d.id = 1 AND ad.year = 1
ORDER BY s.enrollment_no;
