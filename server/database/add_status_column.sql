-- Add status column to tables for soft delete
-- status = 1 (active), status = 0 (deleted)

-- Drop procedure if exists
DROP PROCEDURE IF EXISTS AddStatusColumn;

DELIMITER $$

CREATE PROCEDURE AddStatusColumn()
BEGIN
    -- honour_cards
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'honour_cards' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `honour_cards` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- honour_verticals
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'honour_verticals' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `honour_verticals` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- honour_vertical_courses
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'honour_vertical_courses' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `honour_vertical_courses` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- curriculum_courses
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'curriculum_courses' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `curriculum_courses` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- normal_cards
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'normal_cards' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `normal_cards` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- syllabus
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'syllabus' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `syllabus` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- peo_po_mapping
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'peo_po_mapping' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `peo_po_mapping` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- clusters
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'clusters' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `clusters` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;

    -- cluster_departments
    IF NOT EXISTS (SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cms_local' AND TABLE_NAME = 'cluster_departments' AND COLUMN_NAME = 'status') THEN
        ALTER TABLE `cluster_departments` ADD COLUMN `status` TINYINT(1) DEFAULT 1;
    END IF;
END$$

DELIMITER ;

-- Execute the procedure
CALL AddStatusColumn();

-- Drop the procedure
DROP PROCEDURE AddStatusColumn;
