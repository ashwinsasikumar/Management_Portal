# Chrome Installation Required for PDF Generation

## The Issue
PDF generation requires Chrome or Chromium to convert HTML to PDF. Currently, Chrome is not installed on your system.

## Quick Fix - Install Chrome

### Option 1: Direct Download (Fastest)
1. Download Chrome from: https://www.google.com/chrome/
2. Install the application
3. Restart the backend server
4. PDF generation will work immediately

### Option 2: Using Homebrew
```bash
brew install --cask google-chrome
```

Then restart the server:
```bash
# Kill existing server
lsof -ti:8080 | xargs kill -9

# Start server
cd server
go run .
```

## Temporary Workaround - HTML Preview

While Chrome is being installed, you can view the curriculum in HTML format:

### From Browser
Add `?preview=html` to the PDF URL:
```
http://localhost:8080/api/regulation/4/pdf?preview=html
```

This will show the fully formatted curriculum in your browser, which you can:
- Print to PDF using your browser's print function (Cmd+P → Save as PDF)
- Save as HTML file
- Copy content as needed

### Example URLs
- **PDF (requires Chrome)**: `http://localhost:8080/api/regulation/4/pdf`
- **HTML Preview**: `http://localhost:8080/api/regulation/4/pdf?preview=html`

## After Installing Chrome

1. Verify installation:
```bash
# Should show Chrome location
ls -la "/Applications/Google Chrome.app"
```

2. Restart backend server

3. Test PDF generation - it should work now!

## Why Chrome is Needed

The PDF generator uses:
- **chromedp** - Headless Chrome automation
- Chrome's built-in PDF print functionality
- Better CSS/HTML rendering than alternatives
- Native support for modern web standards

## Alternative: Manual Conversion

If you prefer not to install Chrome:

1. Get HTML preview: `http://localhost:8080/api/regulation/4/pdf?preview=html`
2. Open in any browser
3. Print to PDF (Cmd+P on Mac, Ctrl+P on Windows/Linux)
4. Choose "Save as PDF"

This gives you a PDF without installing Chrome on the server!
