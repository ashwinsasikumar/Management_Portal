-- Table to store mapping between courses and teachers (Workload Allocation)
CREATE TABLE IF NOT EXISTS course_allocations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL,
    teacher_id BIGINT UNSIGNED NOT NULL,
    academic_year VARCHAR(50) NOT NULL, -- e.g., '2025-2026'
    semester INT DEFAULT NULL,          -- The semester this allocation belongs to
    section VARCHAR(10) DEFAULT 'A',    -- Handles multiple sections (A, B, C...)
    role ENUM('Primary', 'Assistant') DEFAULT 'Primary', -- Useful for Labs
    status TINYINT(1) DEFAULT '1',      -- Soft delete (1=Active, 0=Deleted)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign Key Constraints
    CONSTRAINT fk_allocation_course FOREIGN KEY (course_id) 
        REFERENCES courses(course_id) ON DELETE CASCADE,
    CONSTRAINT fk_allocation_teacher FOREIGN KEY (teacher_id) 
        REFERENCES teachers(id) ON DELETE CASCADE,

    -- Constraints
    -- Prevent duplicate assignments of the same teacher to the same course/section in the same academic year
    UNIQUE KEY unique_assignment (course_id, teacher_id, academic_year, section)
) ENGINE=InnoDB;
