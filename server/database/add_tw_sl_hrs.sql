-- Add tw_sl_hrs column to courses table
ALTER TABLE courses ADD COLUMN IF NOT EXISTS tw_sl_hrs INT DEFAULT 0;
