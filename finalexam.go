package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func createTable() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	fmt.Println(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Connect to database error", err)
		return
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS customer (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Println("Cannot create table.", err)
		return
	}

	fmt.Println("Successfully create table.")
}

type customers struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func insertCus(g *gin.Context) {
	c := customers{}

	err := g.ShouldBindJSON(&c)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "JSON parsing on insert error" + err.Error()})
		return
	}
	row := db.QueryRow("INSERT INTO customer (name, email, status) values ($1,$2,$3) RETURNING id", c.Name, c.Email, c.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Insertion error!!! " + err.Error()})
		return
	}
	c.ID = id
	g.JSON(http.StatusCreated, gin.H{
		"id":     c.ID,
		"name":   c.Name,
		"email":  c.Email,
		"status": c.Status,
	})

}

func getOneCust(g *gin.Context) {
	id := g.Param("id")

	var localCust customers

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL select error!!! " + err.Error()})
		return
	}

	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Convert id to num error!!!" + err.Error()})
		return
	}

	row := stmt.QueryRow(idnum)

	err = row.Scan(&localCust.ID, &localCust.Name, &localCust.Email, &localCust.Status)
	if err != nil {
		log.Println("Select id = " + id)
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Select id %d error!!!: %s", idnum, err.Error())})
		return
	}

	g.JSON(http.StatusOK, localCust)
}
func getCustomers(g *gin.Context) {
	var customer []customers
	stmt, err := db.Prepare("select id, name, email, status FROM customer")
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL for selection error!!! " + err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Query to select row error!!! " + err.Error()})
		return
	}
	for rows.Next() {
		c := customers{}
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Status)
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}
		customer = append(customer, c)
	}
	g.JSON(http.StatusOK, customer)
}

func updateCust(g *gin.Context) {
	id := g.Param("id")

	//Pre-select.
	var localCust customers

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL select error!!! " + err.Error()})
		return
	}

	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Convert id to num error!!!" + err.Error()})
		return
	}

	row := stmt.QueryRow(idnum)

	err = row.Scan(&localCust.ID, &localCust.Name, &localCust.Email, &localCust.Status)
	if err != nil {
		log.Println("Select id = " + id)
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Select id %d error!!!: %s", idnum, err.Error())})
		return
	}

	//Actual update.
	c := customers{}

	err = g.ShouldBindJSON(&c)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "JSON parsing Error!!! " + err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE customer SET name=$2,email=$3,status=$4 WHERE id=$1")
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL for update error!!! " + err.Error()})
		return
	}

	_, err = stmt.Exec(id, c.Name, c.Email, c.Status)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Execute update error!!! " + err.Error()})
		return
	}

	g.JSON(http.StatusOK, c)
}

func deleteCust(g *gin.Context) {

	id := g.Param("id")
	//Pre-select.
	var localCust customers

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL select error!!! " + err.Error()})
		return
	}

	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Convert id to num error!!!" + err.Error()})
		return
	}

	row := stmt.QueryRow(idnum)

	err = row.Scan(&localCust.ID, &localCust.Name, &localCust.Email, &localCust.Status)
	if err != nil {
		log.Println("Select id = " + id)
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Select id %d error!!!: %s", idnum, err.Error())})
		return
	}
	//Actual delete.
	stmt, err = db.Prepare("DELETE FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL for delete row error!!! " + err.Error()})
		return
	}

	_, err = stmt.Exec(idnum)
	if err != nil {
		log.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"status": "Execute deletion error!!! " + err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func authMiddleware(g *gin.Context) {
	token := g.GetHeader("Authorization")

	if token != "token2019" {
		g.JSON(http.StatusUnauthorized, gin.H{"status": "Error!!! Unauthorization"})
		g.Abort()
		return
	}

	g.Next()

}

func main() {
	createTable()
	r := gin.Default()
	r.Use(authMiddleware)
	r.POST("/customers", insertCus)
	r.GET("/customers/:id", getOneCust)
	r.GET("/customers", getCustomers)
	r.PUT("/customers/:id", updateCust)
	r.DELETE("/customers/:id", deleteCust)
	r.Run(":2019")
	defer db.Close()
}
