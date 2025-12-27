# Curriculum PDF Generation System

## Overview
This system generates comprehensive curriculum PDFs for departments and regulations using HTML-to-PDF conversion. The PDFs follow official syllabus formatting rules with distinct templates for theory and laboratory courses.

## Architecture

### Backend Components

1. **Models** (`server/models/pdf.go`)
   - `RegulationPDF`: Main data structure containing all regulation information
   - `SemesterPDF`: Semester-level course data
   - `CoursePDF`: Complete course details with syllabus, modules, and mappings
   - `SyllabusModelPDF`, `SyllabusTitlePDF`, `SyllabusTopicPDF`: Module structure

2. **Handler** (`server/handlers/pdf_html.go`)
   - `GenerateRegulationPDFHTML`: Main endpoint handler
   - `fetchCompleteRegulationData`: Fetches all data from database
   - `fetchModelsForPDF`: Retrieves course modules/topics hierarchically
   - `generateHTMLPDF`: Renders HTML and converts to PDF

3. **Route** (`server/routes/routes.go`)
   - Endpoint: `GET /api/regulation/{id}/pdf`

### Frontend Integration

PDF download buttons are integrated into:
- Curriculum Main Page (`curriculumMainPage.js`)
- Regulations Page (`regulationsPage.js`)

## Template Structure

### Header Section (CSE-R2022-V1.x Template)
The PDF follows the official template structure for:
- Cover page with regulation name and academic year
- Vision and Mission statements
- Program Educational Objectives (PEOs)
- Program Outcomes (POs)
- Program Specific Outcomes (PSOs)
- PEO-PO Mapping matrix
- Summary of Credit Distribution (semester-wise course tables)

### Course Description Section

After the header section, courses are rendered dynamically based on type:

#### Theory Courses (26MA101 Template)
- Course code, name, and basic info (L-T-P-C, hours, marks)
- Prerequisites
- Course Objectives
- Course Outcomes
- CO-PO Mapping matrix
- CO-PSO Mapping matrix
- Course Content organized as Modules with titles, hours, and topics
- Teamwork activities with hours
- Self Learning topics with hours
- Text Books and References

#### Laboratory Courses (26PH108 Template)
- Course code, name, and basic info
- Prerequisites
- Course Objectives
- Course Outcomes
- CO-PO and CO-PSO Mappings
- List of Experiments (instead of modules)
- Teamwork activities
- Self Learning topics
- References/Manuals

### Data Source

**All content is fetched from the database:**
- Department overview (vision, mission, PEOs, POs, PSOs)
- Semester and course structure
- Course details (credits, hours, marks)
- Syllabus data from normalized tables:
  - `course_objectives`
  - `course_outcomes`
  - `course_prerequisites`
  - `course_references`
  - `course_teamwork` and `course_teamwork_activities`
  - `course_selflearning`, `course_selflearning_main`, `course_selflearning_internal`
  - `syllabus_models`, `syllabus_titles`, `syllabus_topics`
- Mapping matrices:
  - PEO-PO mapping from `peo_po_mapping`
  - CO-PO mapping from `co_po_mapping`
  - CO-PSO mapping from `co_pso_mapping`

**Reference PDFs are used ONLY for layout guidance**, not as data sources.

## Course Ordering

Courses appear in the exact order defined in the database:
- Ordering is based on the `id` column in the `curriculum_courses` junction table
- This preserves the insertion order and mirrors the reference PDF structure
- Electives and all course types maintain their database-defined positions

## Installation Requirements

### Install Chrome/Chromium (for PDF Generation)

The system uses `chromedp` (headless Chrome) for HTML-to-PDF conversion. Chrome or Chromium must be installed on the system.

#### macOS
Chrome is typically already installed. If not:
```bash
brew install --cask google-chrome
# OR for Chromium:
brew install --cask chromium
```

#### Ubuntu/Debian
```bash
# Install Chromium
sudo apt-get update
sudo apt-get install chromium-browser
```

#### CentOS/RHEL
```bash
sudo yum install chromium
```

#### Windows
Download from: https://www.google.com/chrome/

### Go Dependencies
The required `chromedp` package is already included in `go.mod`. Run:
```bash
cd server
go mod download
```

### Verify Installation
Chrome/Chromium should be in your system PATH. Test with:
```bash
# macOS/Linux
which google-chrome
# OR
which chromium
```

## Usage

### API Endpoint
```
GET http://localhost:8080/api/regulation/{id}/pdf
```

### Frontend
Click the PDF download button (document icon) on any regulation card in:
- Curriculum Main Page
- Regulations Page

### Response
- Content-Type: `application/pdf`
- Downloads with filename: `Regulation_{name}_{year}.pdf`

## Template Customization

The HTML template is defined in `htmlTemplate` constant in `pdf_html.go`. It includes:

### CSS Styling
- A4 page size with proper margins
- Print-optimized styles
- Responsive tables with borders
- Page break controls
- Grid layouts for course information

### Helper Functions
- `add`: Addition for template calculations
- `sub`: Subtraction
- `totalHours`: Calculate total L+T+P hours
- `iterate`: Generate range for loops
- `isTheory`: Detect theory courses
- `isLab`: Detect laboratory courses

## Error Handling

Common issues and solutions:

1. **"Chrome not found" or "exec: not found"**
   - Install Chrome or Chromium as per instructions above
   - Ensure Chrome is in system PATH
   - On Linux, you may need to install additional dependencies

2. **"Failed to fetch regulation data"**
   - Check database connectivity
   - Verify regulation ID exists
   - Check database tables have required data

3. **"PDF conversion failed" or "context deadline exceeded"**
   - Increase timeout in code if needed
   - Check if Chrome has sufficient permissions
   - Verify no Chrome processes are hanging

4. **PDF rendering issues**
   - Check CSS styles in HTML template
   - Verify data is properly formatted
   - Test HTML rendering in browser first

## Database Schema Requirements

Required tables for complete PDF generation:
- `curriculum` (regulation info)
- `department_overview` (vision, mission, PEOs, POs, PSOs)
- `semesters` (semester structure)
- `courses` (course details)
- `curriculum_courses` (junction table for ordering)
- `course_objectives`, `course_outcomes`, `course_prerequisites`, `course_references`
- `course_teamwork`, `course_teamwork_activities`
- `course_selflearning`, `course_selflearning_main`, `course_selflearning_internal`
- `syllabus_models`, `syllabus_titles`, `syllabus_topics`
- `peo_po_mapping`, `co_po_mapping`, `co_pso_mapping`

## Development Notes

### Adding New Sections
1. Update `models/pdf.go` with new data structures
2. Modify `fetchCompleteRegulationData()` to fetch new data
3. Update HTML template to render new sections

### Modifying Styles
Edit the `<style>` block in `htmlTemplate` constant

### Changing Page Layout
Modify CSS `@page` settings and HTML structure

## Performance Considerations

- PDF generation is CPU-intensive
- Each request creates temporary files (automatically cleaned up)
- Large regulations (many courses) may take 10-30 seconds
- Consider implementing:
  - Caching for frequently accessed PDFs
  - Background job processing for large PDFs
  - Progress indicators in frontend

## Future Enhancements

Potential improvements:
1. Add table of contents with page numbers
2. Include course-wise credit summary charts
3. Add watermarks or header/footer with university logo
4. Support for multiple departments in single PDF
5. Custom branding/themes per institution
6. PDF caching with invalidation on data changes
7. Async PDF generation with download queue
8. PDF compression for smaller file sizes
