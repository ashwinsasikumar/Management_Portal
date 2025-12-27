# Database Migration - Add Theory and Activity Hours

## Summary
This migration adds two new columns to the `courses` table:
- `theory_hours` (INT, default 0)
- `activity_hours` (INT, default 0)

## How to Apply

### Option 1: Using MySQL CLI
```bash
mysql -u your_username -p your_database_name < server/db/migrations/add_theory_activity_hours.sql
```

### Option 2: Using MySQL Workbench or phpMyAdmin
1. Open your database connection
2. Navigate to the `courses` table
3. Run the SQL commands from `add_theory_activity_hours.sql`

### Option 3: Manual ALTER TABLE
```sql
USE your_database_name;

ALTER TABLE courses 
ADD COLUMN theory_hours INT DEFAULT 0 AFTER credit,
ADD COLUMN activity_hours INT DEFAULT 0 AFTER theory_hours;
```

## Changes Made

### Backend (Go Server)
1. **Models** ([server/models/curriculum.go](server/models/curriculum.go))
   - Added `TheoryHours` and `ActivityHours` fields to `Course` struct
   - Added `TheoryHours` and `ActivityHours` fields to `CourseWithDetails` struct

2. **Handlers** ([server/handlers/curriculum.go](server/handlers/curriculum.go))
   - Updated `AddCourseToSemester`: Modified INSERT query to include new fields
   - Updated `GetSemesterCourses`: Modified SELECT query to fetch new fields

3. **Edit Handlers** ([server/handlers/curriculum_edit.go](server/handlers/curriculum_edit.go))
   - Updated `UpdateCourse`: Modified UPDATE query to include new fields
   - Added diff tracking for `theory_hours` and `activity_hours`

### Frontend (React)
1. **State Management** ([client/src/pages/semesterDetailPage.js](client/src/pages/semesterDetailPage.js))
   - Added `theory_hours`, `activity_hours`, and `total_hours` to `newCourse` state
   - Added `theory_hours`, `activity_hours`, and `total_hours` to `editCourseData` state
   - Updated form submission handlers to send new fields

2. **Add Course Form**
   - Added "Theory Hrs" input field
   - Added "Activity Hrs" input field
   - Added "Total Hrs" input field

3. **Edit Course Modal**
   - Added "Theory Hrs" input field
   - Added "Activity Hrs" input field  
   - Added "Total Hrs" input field

## Verification

After applying the migration, verify the changes:

```sql
-- Check table structure
DESCRIBE courses;

-- Verify data
SELECT course_code, course_name, theory_hours, activity_hours, lecture_hours, tutorial_hours, total_hours 
FROM courses 
LIMIT 5;
```

## Testing

1. Start the backend server:
   ```bash
   cd server
   go run main.go
   ```

2. Start the frontend:
   ```bash
   cd client
   npm start
   ```

3. Navigate to: http://localhost:3000/regulation/4/curriculum/semester/3

4. Test the new fields:
   - Click "Add Course" and verify all hour fields are visible
   - Fill in the new fields (theory hrs, activity hrs, total hrs)
   - Save and verify the data is stored correctly
   - Edit an existing course and verify the new fields appear
   - Update the values and save

## Rollback (if needed)

If you need to rollback this migration:

```sql
ALTER TABLE courses 
DROP COLUMN theory_hours,
DROP COLUMN activity_hours;
```
