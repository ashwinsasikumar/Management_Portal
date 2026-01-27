package curriculum

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"server/db"
	"server/models"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

// GenerateRegulationPDF handles GET /curriculum/:id/pdf
func GenerateRegulationPDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid curriculum ID", http.StatusBadRequest)
		return
	}

	// Fetch all data for the curriculum
	pdfData, err := fetchRegulationData(curriculumID)
	if err != nil {
		log.Println("Error fetching curriculum data:", err)
		http.Error(w, "Failed to fetch curriculum data", http.StatusInternalServerError)
		return
	}

	// Generate PDF
	pdfBytes, err := generatePDF(pdfData)
	if err != nil {
		log.Println("Error generating PDF:", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Set headers and stream PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=curriculum_%d.pdf", curriculumID))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	w.Write(pdfBytes)
}

func fetchRegulationData(curriculumID int) (*models.RegulationPDF, error) {
	pdfData := &models.RegulationPDF{
		CurriculumID: curriculumID,
	}

	// Fetch curriculum basic info
	err := db.DB.QueryRow("SELECT name, academic_year FROM curriculum WHERE id = ?", curriculumID).
		Scan(&pdfData.RegulationName, &pdfData.AcademicYear)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch curriculum: %v", err)
	}

	// Fetch department overview
	var departmentID int
	err = db.DB.QueryRow(`
		SELECT id, vision 
		FROM curriculum_vision 
		WHERE curriculum_id = ?`, curriculumID).
		Scan(&departmentID, &pdfData.Overview.Vision)

	if err == nil {
		// Fetch mission, PEOs, POs, PSOs from normalized tables (all use curriculum_id FK to curriculum.id)
		pdfData.Overview.Mission = fetchDepartmentListItems(curriculumID, "curriculum_mission", "mission_text")
		pdfData.Overview.PEOs = fetchDepartmentListItems(curriculumID, "curriculum_peos", "peo_text")
		pdfData.Overview.POs = fetchDepartmentListItems(curriculumID, "curriculum_pos", "po_text")
		pdfData.Overview.PSOs = fetchDepartmentListItems(curriculumID, "curriculum_psos", "pso_text")
	}

	// Fetch PEO-PO mapping
	pdfData.PEOPOMapping = make(map[string]int)
	rows, _ := db.DB.Query("SELECT peo_index, po_index, mapping_value FROM peo_po_mapping WHERE curriculum_id = ?", curriculumID)
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
	semRows, err := db.DB.Query("SELECT id, semester_number FROM semesters WHERE regulation_id = ? ORDER BY semester_number", curriculumID)
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

		// Fetch courses for semester
		courseRows, _ := db.DB.Query(`
			SELECT c.course_id, c.course_code, c.course_name, c.course_type, c.category, c.credit,
			       c.lecture_hours, c.tutorial_hours, c.practical_hours, 
			       c.cia_marks, c.see_marks, c.total_marks
			FROM courses c
			INNER JOIN curriculum_courses rc ON c.course_id = rc.course_id
			WHERE rc.curriculum_id = ? AND rc.semester_id = ?
			ORDER BY c.course_code`, curriculumID, semID)

		if courseRows != nil {
			defer courseRows.Close()
			for courseRows.Next() {
				var course models.CoursePDF
				courseRows.Scan(&course.CourseID, &course.CourseCode, &course.CourseName, &course.CourseType,
					&course.Category, &course.Credit, &course.LectureHours, &course.TutorialHours,
					&course.PracticalHours, &course.CIAMarks, &course.SEEMarks, &course.TotalMarks)

				// Fetch syllabus from normalized tables
				course.Syllabus.Objectives, _ = fetchObjectives(course.CourseID)
				course.Syllabus.Outcomes, _ = fetchOutcomes(course.CourseID)
				course.Syllabus.ReferenceList, _ = fetchReferences(course.CourseID)

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
						copsoRows.Scan(&coIdx, &psoIdx, &val)
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

func generatePDF(data *models.RegulationPDF) ([]byte, error) {
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

	// Load template
	tmpl, err := template.New("regulation").Funcs(template.FuncMap{
		"escapeLatex": escapeLatex,
		"add":         func(a, b int) int { return a + b },
		"totalHours":  func(l, t, p int) int { return l + t + p },
		"iterate": func(n int) []int {
			res := make([]int, n)
			for i := 0; i < n; i++ {
				res[i] = i
			}
			return res
		},
	}).Parse(latexTemplate)
	if err != nil {
		return nil, err
	}

	// Render template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, dataWithDate)
	if err != nil {
		return nil, err
	}

	// Write .tex file
	texFile := filepath.Join(tmpDir, "regulation.tex")
	err = ioutil.WriteFile(texFile, buf.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	// Compile with pdflatex
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", tmpDir, texFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("LaTeX compilation error:", string(output))
		return nil, fmt.Errorf("LaTeX compilation failed: %v", err)
	}

	// Run twice for TOC
	cmd = exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", tmpDir, texFile)
	cmd.Run()

	// Read PDF
	pdfFile := filepath.Join(tmpDir, "regulation.pdf")
	pdfBytes, err := ioutil.ReadFile(pdfFile)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}

// Helper function to fetch list items from normalized tables for PDF generation
func fetchDepartmentListItems(departmentID int, tableName, columnName string) []models.DepartmentListItem {
	query := fmt.Sprintf("SELECT id, %s, visibility, source_curriculum_id FROM %s WHERE curriculum_id = ? ORDER BY position", columnName, tableName)
	rows, err := db.DB.Query(query, departmentID)
	if err != nil {
		return []models.DepartmentListItem{}
	}
	defer rows.Close()

	items := []models.DepartmentListItem{}
	for rows.Next() {
		var item models.DepartmentListItem
		var sourceDeptID sql.NullInt64
		if err := rows.Scan(&item.ID, &item.Text, &item.Visibility, &sourceDeptID); err == nil {
			if sourceDeptID.Valid {
				item.SourceCurriculumID = int(sourceDeptID.Int64)
			}
			items = append(items, item)
		}
	}
	return items
}

func escapeLatex(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\textbackslash{}",
		"&", "\\&",
		"%", "\\%",
		"$", "\\$",
		"#", "\\#",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}",
	)
	return replacer.Replace(s)
}

const latexTemplate = `\documentclass[12pt,a4paper]{article}
\usepackage[utf8]{inputenc}
\usepackage[margin=1in]{geometry}
\usepackage{array}
\usepackage{longtable}
\usepackage{booktabs}
\usepackage{fancyhdr}
\usepackage{graphicx}
\usepackage{titlesec}
\usepackage{tocloft}
\usepackage{hyperref}

% Page setup
\pagestyle{fancy}
\fancyhf{}
\fancyhead[L]{\textbf{ {{- escapeLatex .RegulationName -}} }}
\fancyhead[R]{\textbf{Page \thepage}}
\renewcommand{\headrulewidth}{0.5pt}

% Section formatting
\titleformat{\section}{\Large\bfseries}{\thesection}{1em}{}
\titleformat{\subsection}{\large\bfseries}{\thesubsection}{1em}{}

\begin{document}

% Cover Page
\begin{titlepage}
\centering
\vspace*{2cm}
{\Huge\bfseries {{escapeLatex .RegulationName}}\par}
\vspace{1cm}
{\Large {{escapeLatex .AcademicYear}}\par}
\vspace{2cm}
{\Large\bfseries Department Regulation Document\par}
\vfill
{\large Generated on {{.GeneratedDate}}\par}
\end{titlepage}

\tableofcontents
\newpage

% Vision and Mission
\section{Vision and Mission}
\subsection{Vision}
{{escapeLatex .Overview.Vision}}

\subsection{Mission}
\begin{itemize}
{{range .Overview.Mission}}
\item {{escapeLatex .Text}}
{{end}}
\end{itemize}

% PEOs
\section{Program Educational Objectives (PEOs)}
\begin{enumerate}
{{range .Overview.PEOs}}
\item {{escapeLatex .Text}}
{{end}}
\end{enumerate}

% POs
\section{Program Outcomes (POs)}
\begin{enumerate}
{{range .Overview.POs}}
\item {{escapeLatex .Text}}
{{end}}
\end{enumerate}

% PSOs
\section{Program Specific Outcomes (PSOs)}
\begin{enumerate}
{{range .Overview.PSOs}}
\item {{escapeLatex .Text}}
{{end}}
\end{enumerate}

% PEO-PO Mapping
\section{PEO-PO Mapping}
\begin{longtable}{|l|*{12}{c|}}
\hline
\textbf{PEO/PO} {{range $i := iterate 12}}& \textbf{PO{{add $i 1}}} {{end}}\\ \hline
\endfirsthead
\hline
\textbf{PEO/PO} {{range $i := iterate 12}}& \textbf{PO{{add $i 1}}} {{end}}\\ \hline
\endhead
{{range $peoIdx, $peo := .Overview.PEOs}}
\textbf{PEO{{add $peoIdx 1}}} {{range $poIdx := iterate 12}}& {{index $.PEOPOMapping (printf "%d-%d" (add $peoIdx 1) (add $poIdx 1))}} {{end}}\\ \hline
{{end}}
\end{longtable}

% Curriculum
\section{Curriculum}
{{range .Semesters}}
\subsection{Semester {{.Number}}}
\begin{longtable}{|p{2cm}|p{4cm}|c|c|c|c|c|c|c|c|p{2cm}|}
\hline
\textbf{Code} & \textbf{Course} & \textbf{L} & \textbf{T} & \textbf{P} & \textbf{C} & \textbf{Hrs/Wk} & \textbf{CIA} & \textbf{SEE} & \textbf{Total} & \textbf{Category} \\ \hline
\endfirsthead
\hline
\textbf{Code} & \textbf{Course} & \textbf{L} & \textbf{T} & \textbf{P} & \textbf{C} & \textbf{Hrs/Wk} & \textbf{CIA} & \textbf{SEE} & \textbf{Total} & \textbf{Category} \\ \hline
\endhead
{{range .Courses}}
{{escapeLatex .CourseCode}} & {{escapeLatex .CourseName}} & {{.LectureHours}} & {{.TutorialHours}} & {{.PracticalHours}} & {{.Credit}} & {{totalHours .LectureHours .TutorialHours .PracticalHours}} & {{.CIAMarks}} & {{.SEEMarks}} & {{.TotalMarks}} & {{escapeLatex .Category}} \\ \hline
{{end}}
\end{longtable}
{{end}}

% Syllabi
\section{Course Syllabi}
{{range $semIdx, $sem := .Semesters}}
{{range $courseIdx, $course := $sem.Courses}}
\subsection{ {{- escapeLatex $course.CourseCode -}} : {{- escapeLatex $course.CourseName -}} }

\textbf{Course Outcomes:}
\begin{enumerate}
{{range $course.Syllabus.Outcomes}}
\item {{escapeLatex .}}
{{end}}
\end{enumerate}

\textbf{CO-PO Mapping:}
\begin{tabular}{|l|*{12}{c|}}
\hline
\textbf{CO/PO} {{range $i := iterate 12}}& \textbf{PO{{add $i 1}}} {{end}}\\ \hline
{{range $coIdx, $co := $course.Syllabus.Outcomes}}
\textbf{CO{{add $coIdx 1}}} {{range $poIdx := iterate 12}}& {{index $course.COPOMapping (printf "%d-%d" $coIdx (add $poIdx 1))}} {{end}}\\ \hline
{{end}}
\end{tabular}

\textbf{Unit I:} {{escapeLatex $course.Syllabus.Unit1}}

\textbf{Unit II:} {{escapeLatex $course.Syllabus.Unit2}}

\textbf{Unit III:} {{escapeLatex $course.Syllabus.Unit3}}

\textbf{Unit IV:} {{escapeLatex $course.Syllabus.Unit4}}

\textbf{Unit V:} {{escapeLatex $course.Syllabus.Unit5}}

\textbf{Textbooks:}
\begin{enumerate}
{{range $course.Syllabus.Textbooks}}
\item {{escapeLatex .}}
{{end}}
\end{enumerate}

\textbf{References:}
\begin{enumerate}
{{range $course.Syllabus.ReferenceList}}
\item {{escapeLatex .}}
{{end}}
\end{enumerate}

\newpage
{{end}}
{{end}}

\end{document}
`
