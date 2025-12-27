package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server/db"
	"server/models"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
)

// GenerateRegulationPDFHTML handles GET /regulation/:id/pdf using HTML to PDF conversion
func GenerateRegulationPDFHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regulationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid regulation ID", http.StatusBadRequest)
		return
	}

	// Fetch all data for the regulation
	pdfData, err := fetchCompleteRegulationData(regulationID)
	if err != nil {
		log.Println("Error fetching regulation data:", err)
		http.Error(w, "Failed to fetch regulation data", http.StatusInternalServerError)
		return
	}

	// Check if we should return HTML preview (for debugging when Chrome is not installed)
	if r.URL.Query().Get("preview") == "html" {
		generateHTMLPreview(w, pdfData)
		return
	}

	// Generate PDF from HTML
	pdfBytes, err := generateHTMLPDF(pdfData)
	if err != nil {
		log.Println("Error generating PDF:", err)
		// Provide helpful error message about Chrome requirement
		errorMsg := fmt.Sprintf("Failed to generate PDF: %v\n\nChrome/Chromium is required for PDF generation.\n"+
			"Install it with: brew install --cask google-chrome\n\n"+
			"Or download the HTML preview by adding ?preview=html to the URL", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	// Set headers and stream PDF
	w.Header().Set("Content-Type", "application/pdf")
	filename := fmt.Sprintf("Regulation_%s_%s.pdf",
		strings.ReplaceAll(pdfData.RegulationName, " ", "_"),
		strings.ReplaceAll(pdfData.AcademicYear, " ", "_"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	w.Write(pdfBytes)
}

// generateHTMLPreview generates an HTML preview of the curriculum
func generateHTMLPreview(w http.ResponseWriter, data *models.RegulationPDF) {
	type PDFDataWithDate struct {
		*models.RegulationPDF
		GeneratedDate string
	}
	dataWithDate := &PDFDataWithDate{
		RegulationPDF: data,
		GeneratedDate: time.Now().Format("January 2, 2006"),
	}

	tmpl, err := template.New("regulation").Funcs(template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"totalHours": func(l, t, p int) int { return l + t + p },
		"iterate": func(n int) []int {
			result := make([]int, n)
			for i := 0; i < n; i++ {
				result[i] = i
			}
			return result
		},
		"isTheory": func(courseType string) bool {
			ct := strings.ToLower(courseType)
			return strings.Contains(ct, "theory") || (!strings.Contains(ct, "lab") && !strings.Contains(ct, "practical"))
		},
		"isLab": func(courseType string) bool {
			ct := strings.ToLower(courseType)
			return strings.Contains(ct, "lab") || strings.Contains(ct, "practical") || strings.Contains(ct, "experiment")
		},
	}).Parse(htmlTemplate)

	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=curriculum_preview.html"))
	tmpl.Execute(w, dataWithDate)
}

func fetchCompleteRegulationData(regulationID int) (*models.RegulationPDF, error) {
	pdfData := &models.RegulationPDF{
		RegulationID: regulationID,
	}

	// Fetch regulation basic info
	err := db.DB.QueryRow("SELECT name, academic_year FROM curriculum WHERE id = ?", regulationID).
		Scan(&pdfData.RegulationName, &pdfData.AcademicYear)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch curriculum: %v", err)
	}

	// Fetch department overview
	var visionJSON, missionJSON, peosJSON, posJSON, psosJSON []byte
	err = db.DB.QueryRow(`
		SELECT vision, mission, peos, pos, psos 
		FROM department_overview 
		WHERE regulation_id = ?`, regulationID).
		Scan(&visionJSON, &missionJSON, &peosJSON, &posJSON, &psosJSON)

	if err == nil {
		json.Unmarshal(visionJSON, &pdfData.Overview.Vision)
		json.Unmarshal(missionJSON, &pdfData.Overview.Mission)

		// Parse PEOs, POs, PSOs and extract text (they are []string in template usage)
		var peoItems []models.DepartmentListItem
		var poItems []models.DepartmentListItem
		var psoItems []models.DepartmentListItem

		json.Unmarshal(peosJSON, &peoItems)
		json.Unmarshal(posJSON, &poItems)
		json.Unmarshal(psosJSON, &psoItems)

		// Store as DepartmentListItem arrays, template will extract .Text
		pdfData.Overview.PEOs = peoItems
		pdfData.Overview.POs = poItems
		pdfData.Overview.PSOs = psoItems
	}

	// Fetch PEO-PO mapping
	pdfData.PEOPOMapping = make(map[string]int)
	rows, _ := db.DB.Query("SELECT peo_index, po_index, mapping_value FROM peo_po_mapping WHERE regulation_id = ?", regulationID)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var peoIdx, poIdx, val int
			rows.Scan(&peoIdx, &poIdx, &val)
			key := fmt.Sprintf("%d-%d", peoIdx, poIdx)
			pdfData.PEOPOMapping[key] = val
		}
	}

	// Fetch semesters
	semRows, err := db.DB.Query("SELECT id, semester_number FROM semesters WHERE regulation_id = ? ORDER BY semester_number", regulationID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch semesters: %v", err)
	}
	defer semRows.Close()

	for semRows.Next() {
		var semID, semNum int
		semRows.Scan(&semID, &semNum)

		semData := models.SemesterPDF{
			SemesterNumber: semNum,
			Courses:        []models.CoursePDF{},
		}

		// Fetch courses for semester - maintain database order
		courseRows, _ := db.DB.Query(`
			SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.category, c.credit,
			       c.lecture_hours, c.tutorial_hours, c.practical_hours, 
			       c.theory_hours, c.activity_hours, c.total_hours,
			       c.cia_marks, c.see_marks, c.total_marks
			FROM courses c
			INNER JOIN curriculum_courses rc ON c.course_id = rc.course_id
			WHERE rc.regulation_id = ? AND rc.semester_id = ?
			ORDER BY rc.id`, regulationID, semID) // Order by junction table ID to preserve insertion order

		if courseRows != nil {
			defer courseRows.Close()
			for courseRows.Next() {
				var course models.CoursePDF
				courseRows.Scan(&course.CourseID, &course.CourseCode, &course.CourseName, &course.CourseType,
					&course.Category, &course.Credit, &course.LectureHours, &course.TutorialHours,
					&course.PracticalHours, &course.TheoryHours, &course.ActivityHours, &course.TotalHours,
					&course.CIAMarks, &course.SEEMarks, &course.TotalMarks)

				// Fetch syllabus from normalized tables
				course.Syllabus.Objectives, _ = fetchObjectives(course.CourseID)
				course.Syllabus.Outcomes, _ = fetchOutcomes(course.CourseID)
				course.Syllabus.Prerequisites, _ = fetchPrerequisites(course.CourseID)
				course.Syllabus.ReferenceList, _ = fetchReferences(course.CourseID)
				course.Syllabus.Teamwork, _ = fetchTeamwork(course.CourseID)
				course.Syllabus.SelfLearning, _ = fetchSelfLearning(course.CourseID)

				// Fetch models/modules
				course.Models = fetchModelsForPDF(course.CourseID)

				// Fetch CO-PO mapping
				course.COPOMapping = make(map[string]int)
				copoRows, _ := db.DB.Query("SELECT co_index, po_index, mapping_value FROM co_po_mapping WHERE course_id = ?", course.CourseID)
				if copoRows != nil {
					defer copoRows.Close()
					for copoRows.Next() {
						var coIdx, poIdx, val int
						copoRows.Scan(&coIdx, &poIdx, &val)
						key := fmt.Sprintf("%d-%d", coIdx, poIdx)
						course.COPOMapping[key] = val
					}
				}

				// Fetch CO-PSO mapping
				course.COPSOMapping = make(map[string]int)
				copsoRows, _ := db.DB.Query("SELECT co_index, pso_index, mapping_value FROM co_pso_mapping WHERE course_id = ?", course.CourseID)
				if copsoRows != nil {
					defer copsoRows.Close()
					for copsoRows.Next() {
						var coIdx, psoIdx, val int
						copoRows.Scan(&coIdx, &psoIdx, &val)
						key := fmt.Sprintf("%d-%d", coIdx, psoIdx)
						course.COPSOMapping[key] = val
					}
				}

				semData.Courses = append(semData.Courses, course)
			}
		}

		pdfData.Semesters = append(pdfData.Semesters, semData)
	}

	return pdfData, nil
}

func fetchModelsForPDF(courseID int) []models.SyllabusModelPDF {
	modelsList := []models.SyllabusModelPDF{}

	modelRows, err := db.DB.Query(`
		SELECT id, model_name, position 
		FROM syllabus_models 
		WHERE course_id = ? 
		ORDER BY position, id`, courseID)

	if err != nil {
		return modelsList
	}
	defer modelRows.Close()

	for modelRows.Next() {
		var model models.SyllabusModelPDF
		if err := modelRows.Scan(&model.ID, &model.ModelName, &model.Position); err != nil {
			continue
		}

		// Fetch titles for this model
		titleRows, err := db.DB.Query(`
			SELECT id, title_name, hours, position 
			FROM syllabus_titles 
			WHERE model_id = ? 
			ORDER BY position, id`, model.ID)

		if err != nil {
			model.Titles = []models.SyllabusTitlePDF{}
			modelsList = append(modelsList, model)
			continue
		}

		titlesList := []models.SyllabusTitlePDF{}
		for titleRows.Next() {
			var title models.SyllabusTitlePDF
			if err := titleRows.Scan(&title.ID, &title.TitleName, &title.Hours, &title.Position); err != nil {
				continue
			}

			// Fetch topics for this title
			topicRows, err := db.DB.Query(`
				SELECT id, topic, position 
				FROM syllabus_topics 
				WHERE title_id = ? 
				ORDER BY position, id`, title.ID)

			if err != nil {
				title.Topics = []models.SyllabusTopicPDF{}
				titlesList = append(titlesList, title)
				continue
			}

			topicsList := []models.SyllabusTopicPDF{}
			for topicRows.Next() {
				var topic models.SyllabusTopicPDF
				if err := topicRows.Scan(&topic.ID, &topic.Topic, &topic.Position); err != nil {
					continue
				}
				topicsList = append(topicsList, topic)
			}
			topicRows.Close()

			title.Topics = topicsList
			titlesList = append(titlesList, title)
		}
		titleRows.Close()

		model.Titles = titlesList
		modelsList = append(modelsList, model)
	}

	return modelsList
}

func generateHTMLPDF(data *models.RegulationPDF) ([]byte, error) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "regulation_pdf_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Add current date to data
	type PDFDataWithDate struct {
		*models.RegulationPDF
		GeneratedDate string
	}
	dataWithDate := &PDFDataWithDate{
		RegulationPDF: data,
		GeneratedDate: time.Now().Format("January 2, 2006"),
	}

	// Load template with helper functions
	tmpl, err := template.New("regulation").Funcs(template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"totalHours": func(l, t, p int) int { return l + t + p },
		"iterate": func(n int) []int {
			result := make([]int, n)
			for i := 0; i < n; i++ {
				result[i] = i
			}
			return result
		},
		"isTheory": func(courseType string) bool {
			ct := strings.ToLower(courseType)
			return strings.Contains(ct, "theory") || (!strings.Contains(ct, "lab") && !strings.Contains(ct, "practical"))
		},
		"isLab": func(courseType string) bool {
			ct := strings.ToLower(courseType)
			return strings.Contains(ct, "lab") || strings.Contains(ct, "practical") || strings.Contains(ct, "experiment")
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				key := values[i].(string)
				dict[key] = values[i+1]
			}
			return dict
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	// Render template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, dataWithDate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}

	htmlContent := buf.String()

	// Use chromedp to convert HTML to PDF with proper error handling
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var pdfBuf []byte
	err = chromedp.Run(ctx, printToPDF(htmlContent, &pdfBuf))
	if err != nil {
		// Check if the error is due to Chrome not being found
		errMsg := err.Error()
		if strings.Contains(errMsg, "executable file not found") ||
			strings.Contains(errMsg, "chrome") ||
			strings.Contains(errMsg, "Cannot find Chrome") {
			return nil, fmt.Errorf("Chrome/Chromium not found. Please install: brew install --cask google-chrome")
		}
		return nil, fmt.Errorf("PDF conversion failed: %v", err)
	}

	return pdfBuf, nil
}

func printToPDF(html string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf(`
				document.open();
				document.write(%s);
				document.close();
			`, strconv.Quote(html)), nil))
			return err
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var buf []byte
			err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				buf, _, err = page.PrintToPDF().
					WithPrintBackground(true).
					WithScale(1).
					WithPaperWidth(8.27).   // A4 width in inches
					WithPaperHeight(11.69). // A4 height in inches
					WithMarginTop(0.39).    // 10mm
					WithMarginBottom(0.39).
					WithMarginLeft(0.59). // 15mm
					WithMarginRight(0.59).
					Do(ctx)
				return err
			}))
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{.RegulationName}}</title>
	<style>
		@page {
			size: A4;
			margin: 10mm 15mm;
		}
		
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}
		
		body {
			font-family: 'Times New Roman', Times, serif;
			font-size: 11pt;
			line-height: 1.4;
			color: #000;
		}
		
		.cover-page {
			text-align: center;
			padding-top: 30%;
			page-break-after: always;
		}
		
		.cover-page h1 {
			font-size: 24pt;
			font-weight: bold;
			margin-bottom: 20px;
		}
		
		.cover-page h2 {
			font-size: 18pt;
			margin-bottom: 40px;
		}
		
		.cover-page .date {
			font-size: 12pt;
			margin-top: 60px;
		}
		
		h1 {
			font-size: 16pt;
			font-weight: bold;
			margin: 20px 0 10px 0;
			text-align: center;
		}
		
		h2 {
			font-size: 14pt;
			font-weight: bold;
			margin: 16px 0 8px 0;
		}
		
		h3 {
			font-size: 12pt;
			font-weight: bold;
			margin: 12px 0 6px 0;
		}
		
		p {
			margin: 6px 0;
			text-align: justify;
		}
		
		table {
			width: 100%;
			border-collapse: collapse;
			margin: 10px 0;
			font-size: 10pt;
		}
		
		table th,
		table td {
			border: 1px solid #000;
			padding: 4px 6px;
			text-align: left;
			vertical-align: top;
		}
		
		table th {
			background-color: #f0f0f0;
			font-weight: bold;
			text-align: center;
		}
		
		table td.center {
			text-align: center;
		}
		
		ul, ol {
			margin: 8px 0 8px 30px;
		}
		
		li {
			margin: 4px 0;
		}
		
		.page-break {
			page-break-after: always;
		}
		
		.course-section {
			margin: 20px 0;
			page-break-inside: avoid;
		}
		
		.course-header {
			background-color: #e8e8e8;
			padding: 8px;
			margin: 10px 0 5px 0;
			font-weight: bold;
			font-size: 12pt;
		}
		
		.grid-table {
			display: table;
			width: 100%;
			margin: 10px 0;
		}
		
		.grid-row {
			display: table-row;
		}
		
		.grid-cell {
			display: table-cell;
			padding: 4px 8px;
			border: 1px solid #000;
		}
		
		.grid-header {
			font-weight: bold;
			background-color: #f0f0f0;
			width: 150px;
		}
		
		.module {
			margin: 10px 0;
		}
		
		.module-title {
			font-weight: bold;
			margin: 8px 0 4px 0;
		}
		
		.topic-list {
			margin-left: 20px;
		}
		
		.credit-table {
			margin: 15px 0;
		}
		
		.mapping-table {
			font-size: 9pt;
		}
		
		.mapping-table th,
		.mapping-table td {
			padding: 2px 4px;
		}
	</style>
</head>
<body>

<!-- Cover Page -->
<div class="cover-page">
	<h1>{{.RegulationName}}</h1>
	<h2>{{.AcademicYear}}</h2>
	<div class="date">Generated on {{.GeneratedDate}}</div>
</div>

<!-- Vision and Mission -->
<h1>VISION</h1>
<p>{{.Overview.Vision}}</p>

<h1>MISSION</h1>
<ol>
{{range .Overview.Mission}}
	<li>{{.Text}}</li>
{{end}}
</ol>

<!-- PEOs -->
<h1>PROGRAM EDUCATIONAL OBJECTIVES (PEOs)</h1>
<ol>
{{range .Overview.PEOs}}
	<li>{{.Text}}</li>
{{end}}
</ol>

<!-- POs -->
<h1>PROGRAM OUTCOMES (POs)</h1>
<ol>
{{range .Overview.POs}}
	<li>{{.Text}}</li>
{{end}}
</ol>

<!-- PSOs -->
{{if .Overview.PSOs}}
<h1>PROGRAM SPECIFIC OUTCOMES (PSOs)</h1>
<ol>
{{range .Overview.PSOs}}
	<li>{{.Text}}</li>
{{end}}
</ol>
{{end}}

<div class="page-break"></div>

<!-- PEO-PO Mapping -->
<h1>PEO-PO MAPPING</h1>
<table class="mapping-table">
	<thead>
		<tr>
			<th>PEO/PO</th>
			{{range $i := iterate (len .Overview.POs)}}
			<th>PO{{add $i 1}}</th>
			{{end}}
		</tr>
	</thead>
	<tbody>
		{{range $peoIdx, $peo := .Overview.PEOs}}
		<tr>
			<th>PEO{{add $peoIdx 1}}</th>
			{{range $poIdx := iterate (len $.Overview.POs)}}
			<td class="center">{{index $.PEOPOMapping (printf "%d-%d" (add $peoIdx 1) (add $poIdx 1))}}</td>
			{{end}}
		</tr>
		{{end}}
	</tbody>
</table>

<div class="page-break"></div>

<!-- Summary of Credit Distribution -->
<h1>SUMMARY OF CREDIT DISTRIBUTION</h1>
{{range .Semesters}}
<h2>SEMESTER {{.SemesterNumber}}</h2>
<table class="credit-table">
	<thead>
		<tr>
			<th>S.No</th>
			<th>Course Code</th>
			<th>Course Name</th>
			<th>L</th>
			<th>T</th>
			<th>P</th>
			<th>C</th>
			<th>Hours/Week</th>
			<th>CIA</th>
			<th>SEE</th>
			<th>Total</th>
			<th>Category</th>
		</tr>
	</thead>
	<tbody>
		{{range $idx, $course := .Courses}}
		<tr>
			<td class="center">{{add $idx 1}}</td>
			<td>{{$course.CourseCode}}</td>
			<td>{{$course.CourseName}}</td>
			<td class="center">{{$course.LectureHours}}</td>
			<td class="center">{{$course.TutorialHours}}</td>
			<td class="center">{{$course.PracticalHours}}</td>
			<td class="center">{{$course.Credit}}</td>
			<td class="center">{{totalHours $course.LectureHours $course.TutorialHours $course.PracticalHours}}</td>
			<td class="center">{{$course.CIAMarks}}</td>
			<td class="center">{{$course.SEEMarks}}</td>
			<td class="center">{{$course.TotalMarks}}</td>
			<td>{{$course.Category}}</td>
		</tr>
		{{end}}
	</tbody>
</table>
{{end}}

<div class="page-break"></div>

<!-- Course Details -->
<h1>COURSE DESCRIPTIONS</h1>

{{range $semIdx, $sem := .Semesters}}
{{range $courseIdx, $course := $sem.Courses}}

<div class="course-section">
	<div class="course-header">
		{{$course.CourseCode}} - {{$course.CourseName}}
	</div>
	
	<!-- Course Info Grid -->
	<div class="grid-table">
		<div class="grid-row">
			<div class="grid-cell grid-header">Course Code</div>
			<div class="grid-cell">{{$course.CourseCode}}</div>
			<div class="grid-cell grid-header">Category</div>
			<div class="grid-cell">{{$course.Category}}</div>
		</div>
		<div class="grid-row">
			<div class="grid-cell grid-header">L-T-P-C</div>
			<div class="grid-cell">{{$course.LectureHours}}-{{$course.TutorialHours}}-{{$course.PracticalHours}}-{{$course.Credit}}</div>
			<div class="grid-cell grid-header">Hours/Week</div>
			<div class="grid-cell">{{totalHours $course.LectureHours $course.TutorialHours $course.PracticalHours}}</div>
		</div>
		<div class="grid-row">
			<div class="grid-cell grid-header">CIA Marks</div>
			<div class="grid-cell">{{$course.CIAMarks}}</div>
			<div class="grid-cell grid-header">SEE Marks</div>
			<div class="grid-cell">{{$course.SEEMarks}}</div>
		</div>
	</div>
	
	<!-- Prerequisites -->
	{{if $course.Syllabus.Prerequisites}}
	<h3>Prerequisites</h3>
	<ul>
	{{range $course.Syllabus.Prerequisites}}
		<li>{{.}}</li>
	{{end}}
	</ul>
	{{end}}
	
	<!-- Objectives -->
	{{if $course.Syllabus.Objectives}}
	<h3>Course Objectives</h3>
	<ol>
	{{range $course.Syllabus.Objectives}}
		<li>{{.}}</li>
	{{end}}
	</ol>
	{{end}}
	
	<!-- Course Outcomes -->
	{{if $course.Syllabus.Outcomes}}
	<h3>Course Outcomes</h3>
	<p>Upon successful completion of this course, students will be able to:</p>
	<ol>
	{{range $course.Syllabus.Outcomes}}
		<li>{{.}}</li>
	{{end}}
	</ol>
	{{end}}
	
	<!-- CO-PO Mapping -->
	{{if and $course.Syllabus.Outcomes (gt (len $.Overview.POs) 0)}}
	<h3>CO-PO Mapping</h3>
	<table class="mapping-table">
		<thead>
			<tr>
				<th>CO/PO</th>
				{{range $i := iterate (len $.Overview.POs)}}
				<th>PO{{add $i 1}}</th>
				{{end}}
			</tr>
		</thead>
		<tbody>
			{{range $coIdx, $co := $course.Syllabus.Outcomes}}
			<tr>
				<th>CO{{add $coIdx 1}}</th>
				{{range $poIdx := iterate (len $.Overview.POs)}}
				<td class="center">{{index $course.COPOMapping (printf "%d-%d" (add $coIdx 1) (add $poIdx 1))}}</td>
				{{end}}
			</tr>
			{{end}}
		</tbody>
	</table>
	{{end}}
	
	<!-- CO-PSO Mapping -->
	{{if and $course.Syllabus.Outcomes $.Overview.PSOs (gt (len $.Overview.PSOs) 0)}}
	<h3>CO-PSO Mapping</h3>
	<table class="mapping-table">
		<thead>
			<tr>
				<th>CO/PSO</th>
				{{range $i := iterate (len $.Overview.PSOs)}}
				<th>PSO{{add $i 1}}</th>
				{{end}}
			</tr>
		</thead>
		<tbody>
			{{range $coIdx, $co := $course.Syllabus.Outcomes}}
			<tr>
				<th>CO{{add $coIdx 1}}</th>
				{{range $psoIdx := iterate (len $.Overview.PSOs)}}
				<td class="center">{{index $course.COPSOMapping (printf "%d-%d" (add $coIdx 1) (add $psoIdx 1))}}</td>
				{{end}}
			</tr>
			{{end}}
		</tbody>
	</table>
	{{end}}
	
	<!-- Course Content / Modules -->
	{{if $course.Models}}
	<h3>{{if isLab $course.CourseType}}List of Experiments{{else}}Course Content{{end}}</h3>
	{{range $course.Models}}
	<div class="module">
		<div class="module-title">{{.ModelName}}</div>
		{{range .Titles}}
		<div style="margin-left: 15px;">
			<strong>{{.TitleName}}</strong>{{if gt .Hours 0}} ({{.Hours}} hours){{end}}
			{{if .Topics}}
			<div class="topic-list">
			{{range .Topics}}
				<div>• {{.Topic}}</div>
			{{end}}
			</div>
			{{end}}
		</div>
		{{end}}
	</div>
	{{end}}
	{{end}}
	
	<!-- Teamwork -->
	{{if $course.Syllabus.Teamwork}}
	<h3>Teamwork ({{$course.Syllabus.Teamwork.Hours}} hours)</h3>
	<ul>
	{{range $course.Syllabus.Teamwork.Activities}}
		<li>{{.}}</li>
	{{end}}
	</ul>
	{{end}}
	
	<!-- Self Learning -->
	{{if $course.Syllabus.SelfLearning}}
	<h3>Self Learning ({{$course.Syllabus.SelfLearning.Hours}} hours)</h3>
	{{range $course.Syllabus.SelfLearning.MainInputs}}
	<div style="margin: 8px 0;">
		<strong>{{.Main}}</strong>
		{{if .Internal}}
		<ul style="margin-left: 30px;">
		{{range .Internal}}
			<li>{{.}}</li>
		{{end}}
		</ul>
		{{end}}
	</div>
	{{end}}
	{{end}}
	
	<!-- References -->
	{{if $course.Syllabus.ReferenceList}}
	<h3>{{if isLab $course.CourseType}}References / Manuals{{else}}Text Books and References{{end}}</h3>
	<ol>
	{{range $course.Syllabus.ReferenceList}}
		<li>{{.}}</li>
	{{end}}
	</ol>
	{{end}}
</div>

<div class="page-break"></div>

{{end}}
{{end}}

</body>
</html>
`
