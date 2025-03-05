package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jctorrestone/web-service-mr/internal/model"
)

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
	router.GET("/symptoms", getSymptoms)
	router.GET("/symptoms/search", getSymptomsByDesc)
	router.GET("/vital-signs", getVitalSigns)
	//POST
	router.POST("/diseases", postDiseases)
	router.POST("/medicines", postMedicines)
	router.POST("/patients", postPatients)
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

func getPatients(c *gin.Context) {
	var patients []model.Patient
	rows, err := db.Query("SELECT * FROM patient ORDER BY last_name ASC")

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

	c.IndentedJSON(http.StatusOK, patients)
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
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"

	rows, err := db.Query(
		`SELECT * FROM patient 
		WHERE name LIKE ? OR last_name LIKE ? 
		ORDER BY last_name ASC`, query, query)

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

	c.IndentedJSON(http.StatusOK, patients)
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
	rows, err := db.Query("SELECT * FROM disease")

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

	c.IndentedJSON(http.StatusOK, diseases)
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
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	rows, err := db.Query("SELECT * FROM disease WHERE description LIKE ?", query)

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

	c.IndentedJSON(http.StatusOK, diseases)
}

func getSymptoms(c *gin.Context) {
	var symptoms []model.Symptom
	rows, err := db.Query("SELECT * FROM symptom")

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

	c.IndentedJSON(http.StatusOK, symptoms)
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
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
	rows, err := db.Query("SELECT * FROM symptom WHERE description LIKE ?", query)

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

	c.IndentedJSON(http.StatusOK, symptoms)
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
	rows, err := db.Query(
		`SELECT m.id, f.id, s.id, s.description, u.id, u.symbol, u.description, m.name, m.dose 
		FROM medicine AS m 
		INNER JOIN formulation AS f 
		ON formulation_id = f.id 
		INNER JOIN shape AS s 
		ON shape_id = s.id 
		INNER JOIN unit AS u 
		ON unit_id = u.id 
		ORDER BY m.name ASC`)

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

	c.IndentedJSON(http.StatusOK, medicines)
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
	query := c.DefaultQuery("q", "")
	query = "%" + query + "%"
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
		ORDER BY m.name ASC`, query)

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

	c.IndentedJSON(http.StatusOK, medicines)
}
