-- Migration script to rename regulation_id to curriculum_id and other column changes
-- Run this in your MySQL database

-- 1. Rename regulation_id to curriculum_id in curriculum tables
ALTER TABLE curriculum_courses CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE curriculum_vision CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE curriculum_mission CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE curriculum_pos CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE curriculum_peos CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE curriculum_psos CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE peo_po_mapping CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

-- 2. Add source_curriculum_id to curriculum tables if not exists
ALTER TABLE curriculum_mission ADD COLUMN IF NOT EXISTS source_curriculum_id INT DEFAULT NULL;

ALTER TABLE curriculum_peos ADD COLUMN IF NOT EXISTS source_curriculum_id INT DEFAULT NULL;

ALTER TABLE curriculum_pos ADD COLUMN IF NOT EXISTS source_curriculum_id INT DEFAULT NULL;

ALTER TABLE curriculum_psos ADD COLUMN IF NOT EXISTS source_curriculum_id INT DEFAULT NULL;

-- 3. Rename columns in honour_cards table
ALTER TABLE honour_cards CHANGE COLUMN regulation_id curriculum_id INT NOT NULL;

ALTER TABLE honour_cards CHANGE COLUMN semester_number number INT NULL;

ALTER TABLE honour_cards CHANGE COLUMN source_department_id source_curriculum_id INT DEFAULT NULL;

-- 4. Rename source_department_id in normal_cards
ALTER TABLE normal_cards CHANGE COLUMN source_department_id source_curriculum_id INT DEFAULT NULL;

-- 5. Rename columns in sharing_tracking table
ALTER TABLE sharing_tracking CHANGE COLUMN source_dept_id source_curriculum_id INT NOT NULL;

ALTER TABLE sharing_tracking CHANGE COLUMN target_dept_id target_curriculum_id INT NOT NULL;

-- 6. Verify the changes
SHOW COLUMNS FROM curriculum_vision;
SHOW COLUMNS FROM curriculum_mission;
SHOW COLUMNS FROM curriculum_peos;
SHOW COLUMNS FROM curriculum_pos;
SHOW COLUMNS FROM curriculum_psos;
SHOW COLUMNS FROM honour_cards;
SHOW COLUMNS FROM normal_cards;
SHOW COLUMNS FROM sharing_tracking;
