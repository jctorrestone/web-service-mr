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
	router.GET("/patients", getPatients)
	router.GET("/patients/:id", getPatientById)
	//POST
	router.POST("/patients", postPatients)

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

		if err := rows.Scan(&patient.ID, &patient.Name, &patient.Lastname, &patient.Gender); err != nil {
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

	if err := row.Scan(&patient.ID, &patient.Name, &patient.Lastname, &patient.Gender); err != nil {

		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no such patient"})
			return
		}

		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, patient)
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
