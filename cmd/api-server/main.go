package main

import (
	"database/sql"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jctorrestone/web-service-mr/internal/model"
)

const N = 10

var db *sql.DB

func main() {
	connect()
	router := gin.Default()
	//GET
	router.GET("/diseases", getDiseases)
	router.GET("/diseases/search", getDiseasesByDesc)
	router.GET("/exams", getExams)
	router.GET("/formulations", getFormulations)
	router.GET("/medicines", getMedicines)
	router.GET("/medicines/search", getMedicinesByDesc)
	router.GET("/patients", getPatients)
	router.GET("/patients/:id", getPatientById)
	router.GET("/patients/search", getPatientsByName)
	router.GET("/records", getRecords)
	router.GET("/records/:id", getRecordsById)
	router.GET("/records/search", getRecordsByPatient)
	router.GET("/symptoms", getSymptoms)
	router.GET("/symptoms/search", getSymptomsByDesc)
	router.GET("/vital-signs", getVitalSigns)
	//POST
	router.POST("/diseases", postDiseases)
	router.POST("/medicines", postMedicines)
	router.POST("/patients", postPatients)
	router.POST("/records", postRecords)
	router.POST("/symptoms", postSymptoms)

	router.Run("localhost:8080")
}

func connect() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "medical_records",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")
}

func getRecordsNum(sql_count string, args ...any) int64 {
	var total int64

	var row *sql.Row
	if len(args) != 0 {
		row = db.QueryRow(sql_count, args...)
	} else {
		row = db.QueryRow(sql_count)
	}

	if err := row.Scan(&total); err != nil {
		return -1
	}

	return total
}

func getPaginationResponse(sql_count string, page int, args ...any) model.Response {
	var response model.Response

	response.Page = page
	response.PrevPage = -1
	response.NextPage = -1
	response.Total = getRecordsNum(sql_count, args...)
	response.LastPage = int(math.Ceil(float64(response.Total)/N) - 1)

	if response.Page < 0 {
		response.Page = 0
	} else if response.Page > response.LastPage {
		response.Page = response.LastPage
	}

	if response.Page > 0 {
		response.PrevPage = response.Page - 1
	}

	if response.Page < response.LastPage {
		response.NextPage = response.Page + 1
	}

	return response
}

func getPatients(c *gin.Context) {
	var patients []model.Patient
	sql_count := "SELECT COUNT(id) AS total FROM patient"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page)

	rows, err := db.Query("SELECT * FROM patient ORDER BY last_name ASC LIMIT ?, ?", response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var patient model.Patient

		if err := rows.Scan(
			&patient.ID, &patient.Name,
			&patient.Lastname, &patient.Gender); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = patients
	c.IndentedJSON(http.StatusOK, response)
}

func getPatientById(c *gin.Context) {
	var patient model.Patient
	id := c.Param("id")
	row := db.QueryRow("SELECT * FROM patient WHERE id = ?", id)

	if err := row.Scan(
		&patient.ID, &patient.Name,
		&patient.Lastname, &patient.Gender); err != nil {

		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no such patient"})
			return
		}

		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, patient)
}

func getPatientsByName(c *gin.Context) {
	var patients []model.Patient

	sql_count := "SELECT COUNT(id) AS total FROM patient WHERE name LIKE ? OR last_name LIKE ?"
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page, query, query)

	rows, err := db.Query(
		`SELECT * FROM patient 
		WHERE name LIKE ? OR last_name LIKE ? 
		ORDER BY last_name ASC 
		LIMIT ?, ?`, query, query, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var patient model.Patient

		if err := rows.Scan(
			&patient.ID, &patient.Name,
			&patient.Lastname, &patient.Gender); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = patients
	c.IndentedJSON(http.StatusOK, response)
}

func postPatients(c *gin.Context) {
	var patient model.Patient

	if err := c.BindJSON(&patient); err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO patient (name, last_name, gender) VALUES (?, ?, ?)",
		patient.Name, patient.Lastname, patient.Gender)

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	patient.ID = id
	c.IndentedJSON(http.StatusCreated, patient)
}

func getExams(c *gin.Context) {
	var exams []model.Exam
	rows, err := db.Query("SELECT * FROM exam")

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var exam model.Exam

		if err := rows.Scan(&exam.ID, &exam.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		exams = append(exams, exam)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, exams)
}

func getDiseases(c *gin.Context) {
	var diseases []model.Disease

	sql_count := "SELECT COUNT(id) AS total FROM disease"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page)

	rows, err := db.Query("SELECT * FROM disease ORDER BY description ASC LIMIT ?, ?", response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var disease model.Disease

		if err := rows.Scan(&disease.ID, &disease.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		diseases = append(diseases, disease)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = diseases
	c.IndentedJSON(http.StatusOK, response)
}

func postDiseases(c *gin.Context) {
	var disease model.Disease

	if err := c.BindJSON(&disease); err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO disease (description) VALUES (?)",
		disease.Description)

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	disease.ID = id
	c.IndentedJSON(http.StatusCreated, disease)
}

func getDiseasesByDesc(c *gin.Context) {
	var diseases []model.Disease

	sql_count := "SELECT COUNT(id) AS total FROM disease WHERE description LIKE ?"
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page, query)

	rows, err := db.Query("SELECT * FROM disease WHERE description LIKE ? ORDER BY description ASC LIMIT ?, ?", query, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var disease model.Disease

		if err := rows.Scan(&disease.ID, &disease.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		diseases = append(diseases, disease)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = diseases
	c.IndentedJSON(http.StatusOK, response)
}

func getSymptoms(c *gin.Context) {
	var symptoms []model.Symptom

	sql_count := "SELECT COUNT(id) AS total FROM symptom"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page)

	rows, err := db.Query("SELECT * FROM symptom ORDER BY description ASC LIMIT ?, ?", response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var symptom model.Symptom

		if err := rows.Scan(&symptom.ID, &symptom.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		symptoms = append(symptoms, symptom)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = symptoms
	c.IndentedJSON(http.StatusOK, response)
}

func postSymptoms(c *gin.Context) {
	var symptom model.Symptom

	if err := c.BindJSON(&symptom); err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO symptom (description) VALUES (?)",
		symptom.Description)

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	symptom.ID = id
	c.IndentedJSON(http.StatusCreated, symptom)
}

func getSymptomsByDesc(c *gin.Context) {
	var symptoms []model.Symptom

	sql_count := "SELECT COUNT(id) AS total FROM symptom WHERE description LIKE ?"
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page, query)

	rows, err := db.Query("SELECT * FROM symptom WHERE description LIKE ? ORDER BY description ASC LIMIT ?, ?", query, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var symptom model.Symptom

		if err := rows.Scan(&symptom.ID, &symptom.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		symptoms = append(symptoms, symptom)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = symptoms
	c.IndentedJSON(http.StatusOK, response)
}

func getVitalSigns(c *gin.Context) {
	var vital_signs []model.VitalSign
	rows, err := db.Query(
		`SELECT vs.id, u.id, u.symbol, u.description, vs.description 
		FROM vital_sign AS vs 
		INNER JOIN unit AS u
		ON unit_id = u.id 
		ORDER BY vs.description ASC`)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var vital_sign model.VitalSign

		if err := rows.Scan(
			&vital_sign.ID, &vital_sign.UnitObj.ID,
			&vital_sign.UnitObj.Symbol, &vital_sign.UnitObj.Description,
			&vital_sign.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		vital_signs = append(vital_signs, vital_sign)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, vital_signs)
}

func getFormulations(c *gin.Context) {
	var formulations []model.Formulation
	rows, err := db.Query(
		`SELECT f.id, s.id, s.description, u.id, u.symbol, u.description 
		FROM formulation AS f
		INNER JOIN shape AS s 
		ON shape_id = s.id 
		INNER JOIN unit AS u 
		ON unit_id = u.id 
		ORDER BY s.description ASC`)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var formulation model.Formulation

		if err := rows.Scan(
			&formulation.ID, &formulation.ShapeObj.ID,
			&formulation.ShapeObj.Description,
			&formulation.UnitObj.ID, &formulation.UnitObj.Symbol,
			&formulation.UnitObj.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		formulations = append(formulations, formulation)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, formulations)
}

func getMedicines(c *gin.Context) {
	var medicines []model.Medicine

	sql_count := "SELECT COUNT(id) AS total FROM medicine"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page)

	rows, err := db.Query(
		`SELECT m.id, f.id, s.id, s.description, u.id, u.symbol, u.description, m.name, m.dose 
		FROM medicine AS m 
		INNER JOIN formulation AS f 
		ON formulation_id = f.id 
		INNER JOIN shape AS s 
		ON shape_id = s.id 
		INNER JOIN unit AS u 
		ON unit_id = u.id 
		ORDER BY m.name ASC 
		LIMIT ?, ?`, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var medicine model.Medicine

		if err := rows.Scan(
			&medicine.ID, &medicine.FormulationObj.ID,
			&medicine.FormulationObj.ShapeObj.ID,
			&medicine.FormulationObj.ShapeObj.Description,
			&medicine.FormulationObj.UnitObj.ID,
			&medicine.FormulationObj.UnitObj.Symbol,
			&medicine.FormulationObj.UnitObj.Description,
			&medicine.Name, &medicine.Dose); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		medicines = append(medicines, medicine)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = medicines
	c.IndentedJSON(http.StatusOK, response)
}

func postMedicines(c *gin.Context) {
	var medicine model.Medicine

	if err := c.BindJSON(&medicine); err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO medicine (formulation_id, name, dose) VALUES (?, ?, ?)",
		medicine.FormulationObj.ID, medicine.Name, medicine.Dose)

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	medicine.ID = id
	c.IndentedJSON(http.StatusCreated, medicine)
}

func getMedicinesByDesc(c *gin.Context) {
	var medicines []model.Medicine

	sql_count := "SELECT COUNT(id) AS total FROM medicine WHERE name LIKE ?"
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page, query)

	rows, err := db.Query(
		`SELECT m.id, f.id, s.id, s.description, u.id, u.symbol, u.description, m.name, m.dose 
		FROM medicine AS m 
		INNER JOIN formulation AS f 
		ON formulation_id = f.id 
		INNER JOIN shape AS s 
		ON shape_id = s.id 
		INNER JOIN unit AS u 
		ON unit_id = u.id 
		WHERE m.name LIKE ? 
		ORDER BY m.name ASC 
		LIMIT ?, ?`, query, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var medicine model.Medicine

		if err := rows.Scan(
			&medicine.ID, &medicine.FormulationObj.ID,
			&medicine.FormulationObj.ShapeObj.ID,
			&medicine.FormulationObj.ShapeObj.Description,
			&medicine.FormulationObj.UnitObj.ID,
			&medicine.FormulationObj.UnitObj.Symbol,
			&medicine.FormulationObj.UnitObj.Description,
			&medicine.Name, &medicine.Dose); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		medicines = append(medicines, medicine)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = medicines
	c.IndentedJSON(http.StatusOK, response)
}

func getRecords(c *gin.Context) {
	var records []model.Record

	sql_count := "SELECT COUNT(id) AS total FROM record"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page)

	rows, err := db.Query(
		`SELECT r.id, p.id, p.name, p.last_name, r.rdate, r.duration 
		FROM record AS r 
		INNER JOIN patient AS p
		ON patient_id = p.id 
		ORDER BY r.rdate DESC  
		LIMIT ?, ?`, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var record model.Record

		if err := rows.Scan(
			&record.ID, &record.PatientObj.ID,
			&record.PatientObj.Name, &record.PatientObj.Lastname,
			&record.Date, &record.Duration); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = records
	c.IndentedJSON(http.StatusOK, response)
}

func getRecordsByPatient(c *gin.Context) {
	var records []model.Record

	sql_count := `SELECT COUNT(r.id) AS total 
		FROM record AS r 
		INNER JOIN patient AS p 
		ON patient_id = p.id 
		WHERE p.name LIKE ? OR p.last_name LIKE ?`

	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	response := getPaginationResponse(sql_count, page, query, query)

	rows, err := db.Query(
		`SELECT r.id, p.id, p.name, p.last_name, r.rdate, r.duration 
		FROM record AS r 
		INNER JOIN patient AS p
		ON patient_id = p.id 
		WHERE p.name LIKE ? OR p.last_name LIKE ? 
		ORDER BY r.rdate DESC 
		LIMIT ?, ?`, query, query, response.Page*N, N)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var record model.Record

		if err := rows.Scan(
			&record.ID, &record.PatientObj.ID,
			&record.PatientObj.Name, &record.PatientObj.Lastname,
			&record.Date, &record.Duration); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Data = records
	c.IndentedJSON(http.StatusOK, response)
}

// TODO COMPLETE ENDPOINT
func getRecordsById(c *gin.Context) {
	var record model.Record
	var histories []model.DiseaseHistory
	var symptoms []model.Symptom
	var vitalSigns []model.RecordVitalSign
	var idx []model.Disease
	var exams []model.Exam
	var treatments []model.Treatment

	id := c.Param("id")
	row := db.QueryRow(
		`SELECT * FROM record 
		INNER JOIN patient 
		ON record.patient_id=patient.id 
		WHERE record.id = ?`, id)

	if err := row.Scan(
		&record.ID, &record.PatientObj.ID, &record.Date,
		&record.Age, &record.Weight, &record.Height, &record.Duration,
		&record.PatientObj.ID, &record.PatientObj.Name, &record.PatientObj.Lastname, &record.PatientObj.Gender); err != nil {

		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no such medical record"})
			return
		}

		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	//--------------------------------------
	rows, _ := db.Query(
		`SELECT dh.id, dh.record_id, d.id, d.description, dh.description 
		FROM disease_history AS dh 
		INNER JOIN disease AS d 
		ON dh.disease_id=d.id 
		WHERE dh.record_id=?`, id)

	defer rows.Close()

	for rows.Next() {
		var history model.DiseaseHistory

		if err := rows.Scan(
			&history.ID, &history.RecordID, &history.DiseaseID, &history.DiseaseDesc, &history.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	//--------------------------------------
	rows, _ = db.Query(
		`SELECT * FROM symptom 
		WHERE id IN (
			SELECT symptom_id 
			FROM record_symptom 
			WHERE record_id=?
		)`, id)

	defer rows.Close()

	for rows.Next() {
		var symptom model.Symptom

		if err := rows.Scan(
			&symptom.ID, &symptom.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		symptoms = append(symptoms, symptom)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	//--------------------------------------
	rows, _ = db.Query(
		`SELECT * FROM disease 
		WHERE id IN (
			SELECT disease_id FROM idx 
			WHERE record_id=?
		)`, id)

	defer rows.Close()

	for rows.Next() {
		var disease model.Disease

		if err := rows.Scan(
			&disease.ID, &disease.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		idx = append(idx, disease)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	//--------------------------------------
	rows, _ = db.Query(
		`SELECT * FROM exam 
		WHERE id IN (
			SELECT exam_id FROM record_exam 
			WHERE record_id=?
		)`, id)

	defer rows.Close()

	for rows.Next() {
		var exam model.Exam

		if err := rows.Scan(
			&exam.ID, &exam.Description); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		exams = append(exams, exam)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	//--------------------------------------
	rows, _ = db.Query(
		`SELECT rvs.id, rvs.record_id, vs.id, vs.description, u.id, u.symbol, rvs.value 
		FROM record_vital_sign AS rvs
		INNER JOIN vital_sign AS vs 
		ON rvs.vital_sign_id=vs.id 
		INNER JOIN unit AS u 
		ON vs.unit_id=u.id 
		WHERE rvs.record_id=?`, id)

	defer rows.Close()

	for rows.Next() {
		var vitalSign model.RecordVitalSign

		if err := rows.Scan(
			&vitalSign.ID, &vitalSign.RecordID, &vitalSign.VitalSignID,
			&vitalSign.Description, &vitalSign.UnitID, &vitalSign.Symbol, &vitalSign.Value); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		vitalSigns = append(vitalSigns, vitalSign)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	//--------------------------------------
	rows, _ = db.Query(
		`SELECT t.id, t.record_id, m.id, m.name, m.dose, f.id, s.id, s.description, u.id, u.symbol, t.quantity, t.dosage, t.frequency, t.instructions 
		FROM treatment AS t
		INNER JOIN medicine AS m
		ON t.medicine_id=m.id
		INNER JOIN formulation AS f
		ON m.formulation_id=f.id
		INNER JOIN shape AS s
		ON f.shape_id=s.id
		INNER JOIN unit AS u
		ON f.unit_id=u.id
		WHERE t.record_id=?`, id)

	defer rows.Close()

	for rows.Next() {
		var treatment model.Treatment

		if err := rows.Scan(
			&treatment.ID, &treatment.RecordID, &treatment.MedicineID, &treatment.Name,
			&treatment.Dose, &treatment.FormulationID, &treatment.ShapeID, &treatment.Description,
			&treatment.UnitID, &treatment.Symbol, &treatment.Quantity, &treatment.Dosage,
			&treatment.Frequency, &treatment.Instructions); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		treatments = append(treatments, treatment)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	fullRecord := model.FullRecord{
		RecordObj:       record,
		DiseasesHistory: histories,
		Symptoms:        symptoms,
		VitalSigns:      vitalSigns,
		Diseases:        idx,
		Exams:           exams,
		Treatments:      treatments,
	}

	c.IndentedJSON(http.StatusOK, fullRecord)
}

// TODO COMPLETE ENDPOINT
func postRecords(c *gin.Context) {
	var record model.Record

	//TODO BIND EVERY OBJECT
	if err := c.BindJSON(&record); err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO record (patient_id, rdate, age, weight, height, duration) VALUES (?, ?, ?, ?, ?, ?)",
		record.PatientObj.ID, record.Date, record.Age, record.Weight, record.Height, record.Duration)

	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": err.Error()})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	record.ID = id
	c.IndentedJSON(http.StatusCreated, record)
}
