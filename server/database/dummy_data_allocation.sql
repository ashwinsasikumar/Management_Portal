-- Dummy Data for Student-Teacher Allocation (Based on Dump 2026-01-24)

-- 1. Departments (Ensure they exist)
INSERT INTO departments (id, department_name, status) VALUES 
(1, 'Computer Science and Engineering', 1),
(2, 'Information Technology', 1),
(3, 'Electronics and Communication Engineering', 1)
ON DUPLICATE KEY UPDATE department_name=VALUES(department_name);

-- 2. Teachers (Ensure they exist)
INSERT INTO teachers (id, name, email, phone, dept, desg, status) VALUES 
(1, 'Dr. Alan Turing', 'alan@cse.edu', '9876543210', '1', 'Professor', 1),
(2, 'Prof. Ada Lovelace', 'ada@cse.edu', '9876543211', '1', 'Associate Professor', 1),
(3, 'Dr. Grace Hopper', 'grace@it.edu', '9876543212', '2', 'Assistant Professor', 1),
(4, 'Prof. Claude Shannon', 'claude@ece.edu', '9876543213', '3', 'Professor', 1),
(5, 'Dr. John von Neumann', 'john@ece.edu', '9876543214', '3', 'Associate Professor', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), dept=VALUES(dept);

-- 3. Students (15 new students with department_id)
-- CSE Students (Dept ID 1)
INSERT INTO students (student_id, enrollment_no, student_name, gender, dob, age, department_id, status) VALUES
(10, 'CSE24001', 'Alice Johnson', 'Female', '2005-01-15', 19, 1, 1),
(11, 'CSE24002', 'Bob Smith', 'Male', '2005-02-20', 19, 1, 1),
(12, 'CSE24003', 'Charlie Brown', 'Male', '2005-03-25', 19, 1, 1),
(13, 'CSE24004', 'Diana Prince', 'Female', '2005-04-30', 19, 1, 1),
(14, 'CSE24005', 'Evan Wright', 'Male', '2005-05-05', 19, 1, 1)
ON DUPLICATE KEY UPDATE student_name=VALUES(student_name);

-- IT Students (Dept ID 2)
INSERT INTO students (student_id, enrollment_no, student_name, gender, dob, age, department_id, status) VALUES
(20, 'IT24001', 'Fiona Gallagher', 'Female', '2005-06-10', 19, 2, 1),
(21, 'IT24002', 'George Miller', 'Male', '2005-07-15', 19, 2, 1),
(22, 'IT24003', 'Hannah Abbott', 'Female', '2005-08-20', 19, 2, 1),
(23, 'IT24004', 'Ian Malcolm', 'Male', '2005-09-25', 19, 2, 1),
(24, 'IT24005', 'Julia Roberts', 'Female', '2005-10-30', 19, 2, 1)
ON DUPLICATE KEY UPDATE student_name=VALUES(student_name);

-- ECE Students (Dept ID 3)
INSERT INTO students (student_id, enrollment_no, student_name, gender, dob, age, department_id, status) VALUES
(30, 'ECE24001', 'Kevin Hart', 'Male', '2005-11-05', 19, 3, 1),
(31, 'ECE24002', 'Laura Croft', 'Female', '2005-12-10', 19, 3, 1),
(32, 'ECE24003', 'Mike Ross', 'Male', '2006-01-15', 18, 3, 1),
(33, 'ECE24004', 'Nina Simone', 'Female', '2006-02-20', 18, 3, 1),
(34, 'ECE24005', 'Oscar Wilde', 'Male', '2006-03-25', 18, 3, 1)
ON DUPLICATE KEY UPDATE student_name=VALUES(student_name);

-- 4. Academic Details for Students (Sync with legacy text fields)
INSERT INTO academic_details (student_id, batch, year, semester, department, curriculum_id) VALUES
(10, '2024-2028', 2024, 1, 'Computer Science and Engineering', 1),
(11, '2024-2028', 2024, 1, 'Computer Science and Engineering', 1),
(12, '2024-2028', 2024, 1, 'Computer Science and Engineering', 1),
(13, '2024-2028', 2024, 1, 'Computer Science and Engineering', 1),
(14, '2024-2028', 2024, 1, 'Computer Science and Engineering', 1),
(20, '2024-2028', 2024, 1, 'Information Technology', 1),
(21, '2024-2028', 2024, 1, 'Information Technology', 1),
(22, '2024-2028', 2024, 1, 'Information Technology', 1),
(23, '2024-2028', 2024, 1, 'Information Technology', 1),
(24, '2024-2028', 2024, 1, 'Information Technology', 1),
(30, '2024-2028', 2024, 1, 'Electronics and Communication Engineering', 1),
(31, '2024-2028', 2024, 1, 'Electronics and Communication Engineering', 1),
(32, '2024-2028', 2024, 1, 'Electronics and Communication Engineering', 1),
(33, '2024-2028', 2024, 1, 'Electronics and Communication Engineering', 1),
(34, '2024-2028', 2024, 1, 'Electronics and Communication Engineering', 1)
ON DUPLICATE KEY UPDATE batch=VALUES(batch);

-- 5. Dummy Courses
INSERT INTO courses (course_id, course_code, course_name, course_type, credit, status) VALUES
(1, 'CS101', 'Programming in C', 'Theory', 3, 1),
(2, 'CS102', 'Data Structures', 'Theory', 4, 1),
(3, 'IT101', 'Web Technologies', 'Theory', 3, 1),
(4, 'EC101', 'Digital Electronics', 'Theory', 3, 1),
(5, 'CS103L', 'Programming Lab', 'Practical', 2, 1)
ON DUPLICATE KEY UPDATE course_name=VALUES(course_name);

-- 6. Course Allocations
INSERT INTO course_allocations (course_id, teacher_id, academic_year, semester, section, role) VALUES
(1, 1, '2025-2026', 1, 'A', 'Primary'),
(2, 2, '2025-2026', 3, 'A', 'Primary'),
(3, 3, '2025-2026', 1, 'A', 'Primary'),
(4, 4, '2025-2026', 1, 'A', 'Primary'),
(5, 1, '2025-2026', 1, 'A', 'Primary'),
(5, 3, '2025-2026', 1, 'A', 'Assistant')
ON DUPLICATE KEY UPDATE role=VALUES(role);
