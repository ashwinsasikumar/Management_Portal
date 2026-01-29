-- Diagnostic queries to check student data

-- 1. Check if students were inserted
SELECT COUNT(*) as student_count FROM students WHERE student_id BETWEEN 1001 AND 3999;

-- 2. Check if academic_details were inserted
SELECT COUNT(*) as academic_count FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;

-- 3. Check students with their academic details
SELECT 
    s.student_id,
    s.enrollment_no,
    s.student_name,
    ad.department,
    ad.year
FROM students s
LEFT JOIN academic_details ad ON s.student_id = ad.student_id
WHERE s.student_id BETWEEN 1001 AND 3999
LIMIT 10;

-- 4. Check department names in departments table
SELECT id, department_name FROM departments ORDER BY id;

-- 5. Check if department names match exactly
SELECT 
    DISTINCT ad.department as academic_dept,
    d.department_name as dept_name,
    d.id as dept_id
FROM academic_details ad
LEFT JOIN departments d ON ad.department = d.department_name
WHERE ad.student_id BETWEEN 1001 AND 3999;

-- 6. Try the exact query from mapping.go for CSE department (id=1), year=1
SELECT 
    s.student_id,
    COALESCE(s.enrollment_no, '') as enrollment,
    s.student_name,
    COALESCE(ad.department, '') as dept,
    COALESCE(ad.year, 0) as year
FROM students s
INNER JOIN academic_details ad ON s.student_id = ad.student_id
INNER JOIN departments d ON ad.department = d.department_name
WHERE d.id = 1 AND ad.year = 1 AND s.status = 1
ORDER BY s.enrollment_no;
