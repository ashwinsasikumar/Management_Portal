-- Migration: Add theory_hours and activity_hours columns to courses table
-- Date: 2025-12-22

-- Add new columns to the courses table
ALTER TABLE courses 
ADD COLUMN theory_hours INT DEFAULT 0 AFTER credit,
ADD COLUMN activity_hours INT DEFAULT 0 AFTER theory_hours;

-- Update existing records if needed (set to 0 by default)
UPDATE courses SET theory_hours = 0 WHERE theory_hours IS NULL;
UPDATE courses SET activity_hours = 0 WHERE activity_hours IS NULL;

-- Verify the changes
-- SELECT course_code, theory_hours, activity_hours FROM courses LIMIT 5;
