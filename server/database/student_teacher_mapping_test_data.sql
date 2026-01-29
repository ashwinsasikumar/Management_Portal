-- Test Data for Student-Teacher Mapping Feature

-- 1. Ensure curriculum exists (needed for foreign key in academic_details)
INSERT INTO curriculum (id, name, academic_year, status) VALUES 
(290, 'Default Curriculum', '2024-2025', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 2. Ensure departments exist
INSERT INTO departments (id, department_name, status) VALUES 
(1, 'Computer Science and Engineering', 1),
(2, 'Information Technology', 1),
(3, 'Electronics and Communication Engineering', 1),
(4, 'Mechanical Engineering', 1)
ON DUPLICATE KEY UPDATE department_name=VALUES(department_name), status=VALUES(status);

-- 3. Ensure teachers exist
INSERT INTO teachers (id, name, email, phone, dept, desg, status) VALUES 
(101, 'Dr. Alan Turing', 'alan.turing@college.edu', '9876543210', 1, 'Professor', 1),
(102, 'Prof. Ada Lovelace', 'ada.lovelace@college.edu', '9876543211', 1, 'Associate Professor', 1),
(103, 'Dr. Grace Hopper', 'grace.hopper@college.edu', '9876543212', 2, 'Assistant Professor', 1),
(104, 'Prof. Claude Shannon', 'claude.shannon@college.edu', '9876543213', 3, 'Professor', 1),
(105, 'Dr. Donald Knuth', 'donald.knuth@college.edu', '9876543214', 1, 'Professor', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), email=VALUES(email), dept=VALUES(dept), desg=VALUES(desg);

-- 4. Link teachers to departments (department_teachers table)
INSERT INTO department_teachers (teacher_id, department_id, status) VALUES
(101, 1, 1), -- Alan Turing -> CSE
(102, 1, 1), -- Ada Lovelace -> CSE
(103, 2, 1), -- Grace Hopper -> IT
(104, 3, 1), -- Claude Shannon -> ECE
(105, 1, 1)  -- Donald Knuth -> CSE
ON DUPLICATE KEY UPDATE status=VALUES(status);

-- 5. Delete existing test students if they exist (cleanup)
DELETE FROM student_teacher_mapping WHERE student_id BETWEEN 1001 AND 3999;
DELETE FROM academic_details WHERE student_id BETWEEN 1001 AND 3999;
DELETE FROM students WHERE student_id BETWEEN 1001 AND 3999;

-- 6. Create test students with proper IDs
INSERT INTO students (student_id, enrollment_no, register_no, student_name, gender, dob, age, status) VALUES
-- CSE Students (Year 1)
(1001, 'CSE2024001', 'REG2024001', 'Alice Johnson', 'Female', '2005-01-15', 19, 1),
(1002, 'CSE2024002', 'REG2024002', 'Bob Smith', 'Male', '2005-02-20', 19, 1),
(1003, 'CSE2024003', 'REG2024003', 'Charlie Brown', 'Male', '2005-03-25', 19, 1),
(1004, 'CSE2024004', 'REG2024004', 'Diana Prince', 'Female', '2005-04-30', 19, 1),
(1005, 'CSE2024005', 'REG2024005', 'Evan Wright', 'Male', '2005-05-05', 19, 1),
(1006, 'CSE2024006', 'REG2024006', 'Fiona Green', 'Female', '2005-06-10', 19, 1),
(1007, 'CSE2024007', 'REG2024007', 'George Wilson', 'Male', '2005-07-15', 19, 1),
(1008, 'CSE2024008', 'REG2024008', 'Hannah Davis', 'Female', '2005-08-20', 19, 1),

-- IT Students (Year 1)
(2001, 'IT2024001', 'REG2024101', 'Ian Malcolm', 'Male', '2005-09-25', 19, 1),
(2002, 'IT2024002', 'REG2024102', 'Julia Roberts', 'Female', '2005-10-30', 19, 1),
(2003, 'IT2024003', 'REG2024103', 'Kevin Hart', 'Male', '2005-11-05', 19, 1),
(2004, 'IT2024004', 'REG2024104', 'Laura Croft', 'Female', '2005-12-10', 19, 1),
(2005, 'IT2024005', 'REG2024105', 'Mike Ross', 'Male', '2006-01-15', 18, 1),

-- ECE Students (Year 1)
(3001, 'ECE2024001', 'REG2024201', 'Nina Simone', 'Female', '2006-02-20', 18, 1),
(3002, 'ECE2024002', 'REG2024202', 'Oscar Wilde', 'Male', '2006-03-25', 18, 1),
(3003, 'ECE2024003', 'REG2024203', 'Paula Abdul', 'Female', '2006-04-30', 18, 1),
(3004, 'ECE2024004', 'REG2024204', 'Quinn Harper', 'Male', '2006-05-05', 18, 1);

-- 7. Add academic details for all students
INSERT INTO academic_details (student_id, batch, year, semester, degree_level, section, department, student_category, curriculum_id) VALUES
-- CSE Students
(1001, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1002, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1003, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1004, '2024-2028', 1, 1, 'UG', 'A', 'Computer Science and Engineering', 'Regular', 290),
(1005, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),
(1006, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),
(1007, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),
(1008, '2024-2028', 1, 1, 'UG', 'B', 'Computer Science and Engineering', 'Regular', 290),

-- IT Students
(2001, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2002, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2003, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2004, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),
(2005, '2024-2028', 1, 1, 'UG', 'A', 'Information Technology', 'Regular', 290),

-- ECE Students
(3001, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290),
(3002, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290),
(3003, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290),
(3004, '2024-2028', 1, 1, 'UG', 'A', 'Electronics and Communication Engineering', 'Regular', 290);

-- 8. (Optional) Create some sample mappings
-- These can be created through the UI using the "Auto-assign" feature, but here are some examples:
INSERT INTO student_teacher_mapping (student_id, teacher_id, department_id, year, academic_year) VALUES
-- CSE mappings (distribute 8 students among 3 teachers)
(1001, 101, 1, 1, '2024-2025'),
(1002, 101, 1, 1, '2024-2025'),
(1003, 101, 1, 1, '2024-2025'),
(1004, 102, 1, 1, '2024-2025'),
(1005, 102, 1, 1, '2024-2025'),
(1006, 102, 1, 1, '2024-2025'),
(1007, 105, 1, 1, '2024-2025'),
(1008, 105, 1, 1, '2024-2025'),

-- IT mappings (distribute 5 students to 1 teacher)
(2001, 103, 2, 1, '2024-2025'),
(2002, 103, 2, 1, '2024-2025'),
(2003, 103, 2, 1, '2024-2025'),
(2004, 103, 2, 1, '2024-2025'),
(2005, 103, 2, 1, '2024-2025'),

-- ECE mappings (distribute 4 students to 1 teacher)
(3001, 104, 3, 1, '2024-2025'),
(3002, 104, 3, 1, '2024-2025'),
(3003, 104, 3, 1, '2024-2025'),
(3004, 104, 3, 1, '2024-2025');

-- Verification queries (run these to check if data is inserted correctly)
-- SELECT COUNT(*) as total_students FROM students WHERE status = 1;
-- SELECT COUNT(*) as total_teachers FROM teachers WHERE status = 1;
-- SELECT d.department_name, COUNT(DISTINCT t.id) as teacher_count 
-- FROM departments d 
-- LEFT JOIN department_teachers dt ON d.id = dt.department_id AND dt.status = 1
-- LEFT JOIN teachers t ON dt.teacher_id = t.id AND t.status = 1
-- GROUP BY d.id, d.department_name;
-- SELECT d.department_name, ad.year, COUNT(s.student_id) as student_count
-- FROM students s
-- INNER JOIN academic_details ad ON s.student_id = ad.student_id
-- INNER JOIN departments d ON ad.department = d.department_name
-- WHERE s.status = 1
-- GROUP BY d.department_name, ad.year;