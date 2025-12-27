# Quick Start - PDF Generation

## Prerequisites
✅ Chrome or Chromium installed (check with `which google-chrome` or `which chromium`)  
✅ Go server dependencies installed (`go mod download`)  
✅ Backend server running on port 8080  
✅ Frontend running on port 3000

## Installation (One-Time Setup)

### Option 1: Auto Install (Recommended)
```bash
cd server
chmod +x install_pdf_dependencies.sh
./install_pdf_dependencies.sh
```

### Option 2: Manual Install

#### macOS
```bash
# Chrome is usually pre-installed
# If not:
brew install --cask google-chrome
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install chromium-browser
```

## Usage

### From Web Interface
1. Open browser: http://localhost:3000
2. Navigate to "Curriculum" page
3. Hover over any regulation card
4. Click the **PDF download icon** (document icon)
5. PDF will download automatically

### From Command Line
```bash
# Download PDF for regulation ID 4
curl -o curriculum.pdf http://localhost:8080/api/regulation/4/pdf
```

### From Code
```javascript
// React/JavaScript
const downloadPDF = async (regulationId) => {
  const response = await fetch(`http://localhost:8080/api/regulation/${regulationId}/pdf`)
  const blob = await response.blob()
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'curriculum.pdf'
  a.click()
}
```

## What Gets Generated

The PDF includes:
- ✅ Cover page
- ✅ Vision and Mission
- ✅ PEOs, POs, PSOs
- ✅ PEO-PO Mapping matrix
- ✅ Semester-wise credit distribution tables
- ✅ Complete course descriptions with:
  - Course info grid
  - Prerequisites
  - Objectives and Outcomes
  - CO-PO and CO-PSO mappings
  - Course modules/experiments
  - Teamwork activities
  - Self-learning topics
  - References

## Troubleshooting

### "Failed to generate PDF"
**Cause**: Chrome not installed or not in PATH  
**Fix**: 
```bash
# macOS
brew install --cask google-chrome

# Linux
sudo apt-get install chromium-browser
```

### "Failed to fetch regulation data"
**Cause**: Database connection issue or missing data  
**Fix**: 
- Check server logs
- Verify regulation ID exists
- Ensure database has department overview data

### PDF is blank or incomplete
**Cause**: Missing data in database  
**Fix**:
- Add vision, mission, PEOs, POs in Department Overview page
- Add courses to semesters
- Fill in course syllabi

### Generation is slow
**Expected**: 5-15 seconds for typical regulations  
**Why**: Chrome needs to render HTML completely before converting to PDF  
**Note**: First generation may take longer as Chrome initializes

## Quick Test

1. **Start servers** (if not running):
```bash
# Terminal 1 - Backend
cd server
go run main.go

# Terminal 2 - Frontend  
cd client
npm start
```

2. **Generate test PDF**:
- Go to http://localhost:3000
- Click on any regulation
- Fill in Department Overview if empty
- Return to main page
- Click PDF icon on regulation card

3. **Verify**:
- PDF downloads
- Opens correctly
- Shows all sections
- Formatting looks good

## Next Steps

- See `IMPLEMENTATION_SUMMARY.md` for complete details
- See `server/PDF_GENERATION_README.md` for technical documentation
- Customize template in `server/handlers/pdf_html.go` if needed

## Support

- Error logs: Check server console output
- Database issues: Verify MySQL connection
- Template issues: Check `server/handlers/pdf_html.go`
- Frontend issues: Check browser console

---

**Ready to use!** 🚀  
The PDF generation system is fully implemented and ready for testing.
