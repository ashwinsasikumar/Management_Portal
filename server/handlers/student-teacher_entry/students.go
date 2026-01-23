package studentteacher

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetStudents retrieves all students from the database
func GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := `
		SELECT 
			student_id, enrollment_no, register_no, dte_reg_no, 
			application_no, admission_no, student_name, gender, dob, age,
			father_name, mother_name, guardian_name, religion, nationality,
			community, mother_tongue, blood_group, aadhar_no, parent_occupation,
			designation, place_of_work, parent_income, status
		FROM students
		ORDER BY student_id DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying students: %v", err)
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.StudentID, &student.EnrollmentNo, &student.RegisterNo,
			&student.DTERegNo, &student.ApplicationNo, &student.AdmissionNo,
			&student.StudentName, &student.Gender, &student.DOB, &student.Age,
			&student.FatherName, &student.MotherName, &student.GuardianName,
			&student.Religion, &student.Nationality, &student.Community,
			&student.MotherTongue, &student.BloodGroup, &student.AadharNo,
			&student.ParentOccupation, &student.Designation, &student.PlaceOfWork,
			&student.ParentIncome, &student.Status,
		)
		if err != nil {
			log.Printf("Error scanning student row: %v", err)
			continue
		}
		students = append(students, student)
	}

	if students == nil {
		students = []models.Student{}
	}

	json.NewEncoder(w).Encode(students)
}

// GetStudent retrieves a single student by ID
func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	studentID := vars["id"]

	query := `
		SELECT 
			student_id, enrollment_no, register_no, dte_reg_no, 
			application_no, admission_no, student_name, gender, dob, age,
			father_name, mother_name, guardian_name, religion, nationality,
			community, mother_tongue, blood_group, aadhar_no, parent_occupation,
			designation, place_of_work, parent_income, status
		FROM students
		WHERE student_id = ?
	`

	var student models.Student
	err := db.DB.QueryRow(query, studentID).Scan(
		&student.StudentID, &student.EnrollmentNo, &student.RegisterNo,
		&student.DTERegNo, &student.ApplicationNo, &student.AdmissionNo,
		&student.StudentName, &student.Gender, &student.DOB, &student.Age,
		&student.FatherName, &student.MotherName, &student.GuardianName,
		&student.Religion, &student.Nationality, &student.Community,
		&student.MotherTongue, &student.BloodGroup, &student.AadharNo,
		&student.ParentOccupation, &student.Designation, &student.PlaceOfWork,
		&student.ParentIncome, &student.Status,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error querying student: %v", err)
		http.Error(w, "Failed to fetch student", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}

// CreateStudent creates a new student record with all related details in transaction
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req models.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	// INSERT into students table
	insertStudentQuery := `
		INSERT INTO students (
			enrollment_no, register_no, dte_reg_no, application_no,
			admission_no, student_name, gender, dob, age, father_name, mother_name,
			guardian_name, religion, nationality, community, mother_tongue,
			blood_group, aadhar_no, parent_occupation, designation, place_of_work,
			parent_income, status
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	result, err := tx.Exec(
		insertStudentQuery,
		req.EnrollmentNo, req.RegisterNo, req.DTERegNo, req.ApplicationNo,
		req.AdmissionNo, req.StudentName, req.Gender, req.DOB, parseInt(req.Age),
		req.FatherName, req.MotherName, req.GuardianName, req.Religion,
		req.Nationality, req.Community, req.MotherTongue, req.BloodGroup,
		req.AadharNo, req.ParentOccupation, req.Designation, req.PlaceOfWork,
		parseFloat(req.ParentIncome), 1,
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting student: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	studentID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Printf("Error getting student ID: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	// INSERT into academic_details if provided
	if req.Department != "" || req.Batch != "" {
		acadQuery := `
			INSERT INTO academic_details (
				student_id, batch, year, semester, degree_level, section, department,
				student_category, branch_type, seat_category, regulation, quota,
				university, year_of_admission, year_of_completion, student_status, curriculum_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := tx.Exec(
			acadQuery,
			studentID, req.Batch, parseInt(req.Year), parseInt(req.Semester), req.DegreeLevel,
			req.Section, req.Department, req.StudentCategory, req.BranchType,
			req.SeatCategory, req.Regulation, req.Quota, req.University,
			parseInt(req.YearOfAdmission), parseInt(req.YearOfCompletion), req.StudentStatus, parseNullableInt(req.CurriculumID),
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting academic details: %v", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}
	}

	// INSERT into address if provided
	if req.PermanentAddress != "" || req.PresentAddress != "" {
		addrQuery := `
			INSERT INTO address (student_id, permanent_address, present_address, residence_location)
			VALUES (?, ?, ?, ?)
		`
		_, err := tx.Exec(addrQuery, studentID, req.PermanentAddress, req.PresentAddress, req.ResidenceLocation)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting address: %v", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}
	}

	// INSERT into admission_payment if provided
	if req.ReceiptNo != "" || req.Amount != "" {
		payQuery := `
			INSERT INTO admission_payment (
				student_id, dte_register_no, dte_admission_no, receipt_no, receipt_date, amount, bank_name
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err := tx.Exec(
			payQuery, studentID, req.DTERegisterNo, req.DTEAdmissionNo,
			req.ReceiptNo, req.ReceiptDate, parseFloat(req.Amount), req.BankName,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting admission payment: %v", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}
	}

	// INSERT into contact_details if provided
	if req.StudentEmail != "" || req.ParentMobile != "" {
		contactQuery := `
			INSERT INTO contact_details (
				student_id, parent_mobile, student_mobile, student_email, parent_email, official_email
			) VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err := tx.Exec(
			contactQuery, studentID, req.ParentMobile, req.StudentMobile, req.StudentEmail,
			req.ParentEmail, req.OfficialEmail,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting contact details: %v", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}
	}

	// INSERT into hostel_details if provided
	if req.HostelName != "" {
		hostelQuery := `
			INSERT INTO hostel_details (
				student_id, hosteller_type, hostel_name, room_no, room_capacity,
				room_type, floor_no, warden_name, alternate_warden, class_advisor, status
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := tx.Exec(
			hostelQuery, studentID, req.HostellerType, req.HostelName, req.RoomNo,
			parseInt(req.RoomCapacity), req.RoomType, parseInt(req.FloorNo), req.WardenName,
			req.AlternateWarden, req.ClassAdvisor, 1,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting hostel details: %v", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}
	}

	// INSERT into insurance_details if provided
	if req.NomineeName != "" {
		insQuery := `
			INSERT INTO insurance_details (student_id, nominee_name, relationship, nominee_age, status)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err := tx.Exec(insQuery, studentID, req.NomineeName, req.Relationship, parseInt(req.NomineeAge), 1)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting insurance details: %v", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}
	}

	// INSERT into school_details for each school if provided
	if len(req.SchoolDetails) > 0 {
		schoolQuery := `
			INSERT INTO school_details (
				student_id, school_name, board, year_of_pass, state, tc_no, tc_date, total_marks, status
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		for _, school := range req.SchoolDetails {
			if school.SchoolName != "" {
				_, err := tx.Exec(
					schoolQuery, studentID, school.SchoolName, school.Board, parseInt(school.YearOfPass),
					school.State, school.TCNo, school.TCDate, parseFloat(school.TotalMarks), 1,
				)
				if err != nil {
					tx.Rollback()
					log.Printf("Error inserting school details: %v", err)
					http.Error(w, "Failed to create student", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	// Return created student
	student := models.Student{
		StudentID:        int(studentID),
		EnrollmentNo:     req.EnrollmentNo,
		RegisterNo:       req.RegisterNo,
		DTERegNo:         req.DTERegNo,
		ApplicationNo:    req.ApplicationNo,
		AdmissionNo:      req.AdmissionNo,
		StudentName:      req.StudentName,
		Gender:           req.Gender,
		DOB:              req.DOB,
		Age:              parseInt(req.Age),
		FatherName:       req.FatherName,
		MotherName:       req.MotherName,
		GuardianName:     req.GuardianName,
		Religion:         req.Religion,
		Nationality:      req.Nationality,
		Community:        req.Community,
		MotherTongue:     req.MotherTongue,
		BloodGroup:       req.BloodGroup,
		AadharNo:         req.AadharNo,
		ParentOccupation: req.ParentOccupation,
		Designation:      req.Designation,
		PlaceOfWork:      req.PlaceOfWork,
		ParentIncome:     parseFloat(req.ParentIncome),
		Status:           1,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// UpdateStudent updates an existing student record and all optional related tables
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	studentID := vars["id"]

	var req models.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// UPDATE students table
	updateQuery := `
		UPDATE students SET
			enrollment_no = ?, register_no = ?, dte_reg_no = ?,
			application_no = ?, admission_no = ?, student_name = ?, gender = ?,
			dob = ?, age = ?, father_name = ?, mother_name = ?,
			guardian_name = ?, religion = ?, nationality = ?, community = ?,
			mother_tongue = ?, blood_group = ?, aadhar_no = ?,
			parent_occupation = ?, designation = ?, place_of_work = ?,
			parent_income = ?
		WHERE student_id = ?
	`

	result, err := tx.Exec(
		updateQuery,
		req.EnrollmentNo, req.RegisterNo, req.DTERegNo, req.ApplicationNo,
		req.AdmissionNo, req.StudentName, req.Gender, req.DOB, parseInt(req.Age),
		req.FatherName, req.MotherName, req.GuardianName, req.Religion,
		req.Nationality, req.Community, req.MotherTongue, req.BloodGroup,
		req.AadharNo, req.ParentOccupation, req.Designation, req.PlaceOfWork,
		parseFloat(req.ParentIncome), studentID,
	)

	if err != nil {
		log.Printf("Error updating student: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// UPDATE academic_details if provided
	if req.Department != "" || req.Batch != "" {
		academicQuery := `
			INSERT INTO academic_details (
				student_id, batch, year, semester, degree_level, section, department,
				student_category, branch_type, seat_category, regulation, quota,
				university, year_of_admission, year_of_completion, student_status, curriculum_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			batch = ?, year = ?, semester = ?, degree_level = ?, section = ?,
			department = ?, student_category = ?, branch_type = ?, seat_category = ?,
			regulation = ?, quota = ?, university = ?, year_of_admission = ?,
			year_of_completion = ?, student_status = ?, curriculum_id = ?
		`
		_, err := tx.Exec(academicQuery, studentID, req.Batch, parseInt(req.Year), parseInt(req.Semester),
			req.DegreeLevel, req.Section, req.Department, req.StudentCategory,
			req.BranchType, req.SeatCategory, req.Regulation, req.Quota,
			req.University, parseInt(req.YearOfAdmission), parseInt(req.YearOfCompletion), req.StudentStatus,
			parseNullableInt(req.CurriculumID), req.Batch, parseInt(req.Year), parseInt(req.Semester), req.DegreeLevel,
			req.Section, req.Department, req.StudentCategory, req.BranchType,
			req.SeatCategory, req.Regulation, req.Quota, req.University,
			parseInt(req.YearOfAdmission), parseInt(req.YearOfCompletion), req.StudentStatus, parseNullableInt(req.CurriculumID))
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating academic details: %v", err)
			http.Error(w, "Failed to update academic details", http.StatusInternalServerError)
			return
		}
	}

	// UPDATE address if provided
	if req.PermanentAddress != "" || req.PresentAddress != "" {
		addressQuery := `
			INSERT INTO address (student_id, permanent_address, present_address, residence_location)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			permanent_address = ?, present_address = ?, residence_location = ?
		`
		_, err := tx.Exec(addressQuery, studentID, req.PermanentAddress, req.PresentAddress,
			req.ResidenceLocation, req.PermanentAddress, req.PresentAddress, req.ResidenceLocation)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating address: %v", err)
			http.Error(w, "Failed to update address", http.StatusInternalServerError)
			return
		}
	}

	// UPDATE admission_payment if provided
	if req.ReceiptNo != "" || req.Amount != "" {
		paymentQuery := `
			INSERT INTO admission_payment (
				student_id, dte_register_no, dte_admission_no, receipt_no, receipt_date, amount, bank_name
			) VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			dte_register_no = ?, dte_admission_no = ?, receipt_no = ?,
			receipt_date = ?, amount = ?, bank_name = ?
		`
		_, err := tx.Exec(paymentQuery, studentID, req.DTERegisterNo, req.DTEAdmissionNo,
			req.ReceiptNo, req.ReceiptDate, parseFloat(req.Amount), req.BankName,
			req.DTERegisterNo, req.DTEAdmissionNo, req.ReceiptNo, req.ReceiptDate, parseFloat(req.Amount), req.BankName)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating admission payment: %v", err)
			http.Error(w, "Failed to update admission payment", http.StatusInternalServerError)
			return
		}
	}

	// UPDATE contact_details if provided
	if req.StudentEmail != "" || req.ParentMobile != "" {
		contactQuery := `
			INSERT INTO contact_details (
				student_id, parent_mobile, student_mobile, student_email, parent_email, official_email
			) VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			parent_mobile = ?, student_mobile = ?, student_email = ?,
			parent_email = ?, official_email = ?
		`
		_, err := tx.Exec(contactQuery, studentID, req.ParentMobile, req.StudentMobile,
			req.StudentEmail, req.ParentEmail, req.OfficialEmail,
			req.ParentMobile, req.StudentMobile, req.StudentEmail, req.ParentEmail, req.OfficialEmail)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating contact details: %v", err)
			http.Error(w, "Failed to update contact details", http.StatusInternalServerError)
			return
		}
	}

	// UPDATE hostel_details if provided
	if req.HostelName != "" {
		hostelQuery := `
			INSERT INTO hostel_details (
				student_id, hosteller_type, hostel_name, room_no, room_capacity, room_type,
				floor_no, warden_name, alternate_warden, class_advisor
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			hosteller_type = ?, hostel_name = ?, room_no = ?, room_capacity = ?,
			room_type = ?, floor_no = ?, warden_name = ?, alternate_warden = ?, class_advisor = ?
		`
		_, err := tx.Exec(hostelQuery, studentID, req.HostellerType, req.HostelName,
			req.RoomNo, parseInt(req.RoomCapacity), req.RoomType, parseInt(req.FloorNo), req.WardenName,
			req.AlternateWarden, req.ClassAdvisor,
			req.HostellerType, req.HostelName, req.RoomNo, parseInt(req.RoomCapacity), req.RoomType,
			parseInt(req.FloorNo), req.WardenName, req.AlternateWarden, req.ClassAdvisor)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating hostel details: %v", err)
			http.Error(w, "Failed to update hostel details", http.StatusInternalServerError)
			return
		}
	}

	// UPDATE insurance_details if provided
	if req.NomineeName != "" {
		insuranceQuery := `
			INSERT INTO insurance_details (student_id, nominee_name, relationship, nominee_age, status)
			VALUES (?, ?, ?, ?, 1)
			ON DUPLICATE KEY UPDATE
			nominee_name = ?, relationship = ?, nominee_age = ?
		`
		_, err := tx.Exec(insuranceQuery, studentID, req.NomineeName, req.Relationship,
			parseInt(req.NomineeAge), req.NomineeName, req.Relationship, parseInt(req.NomineeAge))
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating insurance details: %v", err)
			http.Error(w, "Failed to update insurance details", http.StatusInternalServerError)
			return
		}
	}

	// UPDATE school_details - delete old and insert new
	if len(req.SchoolDetails) > 0 {
		// Delete existing school records
		_, err := tx.Exec(`DELETE FROM school_details WHERE student_id = ?`, studentID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error deleting school details: %v", err)
			http.Error(w, "Failed to update school details", http.StatusInternalServerError)
			return
		}

		// Insert new school records
		schoolQuery := `
			INSERT INTO school_details (
				student_id, school_name, board, year_of_pass, state, tc_no, tc_date, total_marks, status
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		for _, school := range req.SchoolDetails {
			if school.SchoolName != "" {
				_, err := tx.Exec(
					schoolQuery, studentID, school.SchoolName, school.Board, parseInt(school.YearOfPass),
					school.State, school.TCNo, school.TCDate, parseFloat(school.TotalMarks), 1,
				)
				if err != nil {
					tx.Rollback()
					log.Printf("Error inserting school details: %v", err)
					http.Error(w, "Failed to update school details", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	// Return updated student
	updatedStudent := models.Student{
		StudentID:        parseInt(studentID),
		EnrollmentNo:     req.EnrollmentNo,
		RegisterNo:       req.RegisterNo,
		DTERegNo:         req.DTERegNo,
		ApplicationNo:    req.ApplicationNo,
		AdmissionNo:      req.AdmissionNo,
		StudentName:      req.StudentName,
		Gender:           req.Gender,
		DOB:              req.DOB,
		Age:              parseInt(req.Age),
		FatherName:       req.FatherName,
		MotherName:       req.MotherName,
		GuardianName:     req.GuardianName,
		Religion:         req.Religion,
		Nationality:      req.Nationality,
		Community:        req.Community,
		MotherTongue:     req.MotherTongue,
		BloodGroup:       req.BloodGroup,
		AadharNo:         req.AadharNo,
		ParentOccupation: req.ParentOccupation,
		Designation:      req.Designation,
		PlaceOfWork:      req.PlaceOfWork,
		ParentIncome:     parseFloat(req.ParentIncome),
		Status:           1,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStudent)
}

// DeleteStudent deletes a student record
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	studentID := vars["id"]

	query := `DELETE FROM students WHERE student_id = ?`
	result, err := db.DB.Exec(query, studentID)
	if err != nil {
		log.Printf("Error deleting student: %v", err)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student deleted successfully",
	})
}

// Helper functions to parse string values from form input
func parseString(s string) string {
	return s
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	val, _ := strconv.Atoi(s)
	return val
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

// parseNullableInt converts string to *int, returns nil if empty or zero
func parseNullableInt(s string) *int {
	val := parseInt(s)
	if val == 0 {
		return nil
	}
	return &val
}
