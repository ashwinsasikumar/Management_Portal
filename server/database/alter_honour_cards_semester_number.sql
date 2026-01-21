-- Make semester_number nullable in honour_cards table
-- This allows honour cards to be created without a semester number

ALTER TABLE `honour_cards` 
MODIFY COLUMN `semester_number` int NULL;

-- Update any existing records with semester_number = 0 to NULL
UPDATE `honour_cards` 
SET `semester_number` = NULL 
WHERE `semester_number` = 0;
