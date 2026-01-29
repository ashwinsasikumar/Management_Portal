-- Quick fix for student data to make them visible in the UI

-- First, let's check what we have
SELECT 'Checking students...' as step;
SELECT student_id, enrollment_no, student_name, status FROM students WHERE student_id BETWEEN 1001 AND 3999;

SELECT 'Checking academic_details...' as step;
SELECT student_id, department, year FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- Check if department names match exactly
SELECT 'Checking department name matching...' as step;
SELECT 
    ad.student_id,
    ad.department as academic_dept,
    d.department_name as dept_table_name,
    CASE WHEN d.id IS NULL THEN 'NO MATCH' ELSE 'MATCHED' END as match_status
FROM academic_details ad
LEFT JOIN departments d ON ad.department = d.department_name
WHERE ad.student_id BETWEEN 1001 AND 3999;

-- If departments don't match, this will show the actual names
SELECT 'Department names in departments table:' as info;
SELECT id, department_name, CHAR_LENGTH(department_name) as name_length FROM departments;

-- Test the actual query that the UI uses for CSE students
SELECT 'Testing actual UI query for CSE (dept_id=1, year=1):' as step;
SELECT 
    s.student_id,
    s.enrollment_no,
    s.student_name,
    ad.department,
    ad.year,
    d.id as dept_id
FROM students s
INNER JOIN academic_details ad ON s.student_id = ad.student_id
INNER JOIN departments d ON ad.department = d.department_name
WHERE d.id = 1 AND ad.year = 1 AND s.status = 1
ORDER BY s.enrollment_no;
