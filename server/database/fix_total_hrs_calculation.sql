-- Fix total_hrs calculation to include practical_total_hrs
-- This migration updates the GENERATED column formula for total_hrs

-- Drop the existing generated column
ALTER TABLE `courses` DROP COLUMN `total_hrs`;

-- Re-add the column with the correct formula including practical_total_hrs
ALTER TABLE `courses` 
ADD COLUMN `total_hrs` int GENERATED ALWAYS AS (
    (`theory_total_hrs` + `activity_total_hrs` + `tutorial_total_hrs` + COALESCE(`practical_total_hrs`, 0))
) STORED;
