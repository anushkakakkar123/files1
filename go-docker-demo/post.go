package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "some-postgres"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "form"
)

type Employee struct {
	Name      string `json:"name"`
	LeaveType string `json:"leave_type"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Teams     string `json:"teams"`
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetAllEmployees(db *sql.DB) []Employee {
	var employees []Employee
	query := `SELECT * FROM form1.leavestype2`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Failed to get data from database:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee

		err := rows.Scan(&employee.Name, &employee.LeaveType, &employee.StartDate, &employee.EndDate, &employee.Teams)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
			return nil
		}
		employees = append(employees, employee)
	}
	return employees
}

func leave_form(c *gin.Context, db *sql.DB) {
	var data Employee
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	insertQuery := `INSERT INTO form1.leavestype2 (name, leave_type, start_date, end_date, teams) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(insertQuery, data.Name, data.LeaveType, data.StartDate, data.EndDate, data.Teams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Leave request submitted successfully"})
}

func main() {
	router := gin.Default()
	// Connect to the database
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()
	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST"}
	router.Use(cors.New(config))
	// Routes for GET and POST methods
	router.GET("/getData", func(c *gin.Context) {
		leaveRequests := GetAllEmployees(db)
		if leaveRequests != nil {
			c.JSON(http.StatusOK, leaveRequests)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
	})
	router.POST("/postData", func(c *gin.Context) { leave_form(c, db) })

	router.Run(":8080")
}
