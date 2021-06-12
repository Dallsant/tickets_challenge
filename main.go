package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error
var postgressPassword = os.Getenv("PGP")

var (
	outfile, _ = os.OpenFile("./err.log", os.O_RDWR|os.O_APPEND, 0755)
	l          = log.New(outfile, "", 0)
)

var SigningKey = []byte("mySigningKey")

func main() {
	fmt.Println("Starting App at port 8000")
	postgressPassword := "postgres"
	pgurl := fmt.Sprintf("host=localhost port=5434 user=postgres dbname=tickets_challenge password=%s sslmode=disable", postgressPassword)
	db, err = gorm.Open("postgres", pgurl)

	if err != nil {
		logger(err.Error(), "PostgresDB conexion")
	} else if err == nil {
		fmt.Println("Postgress: Connected")
	}

	db.SingularTable(true)
	InitialMigration()

	handleRequests()
}

type Connection struct {
	gorm.Model
	Username string `json:"username"`
}

func InitialMigration() {
	db.AutoMigrate(&Ticket{}, &User{})
}

func logger(err string, domain string) {
	fmt.Println(fmt.Sprintf("$timestamp: %v $at: %s $message: %v", time.Now().Unix(), domain, err))

	l.Println(fmt.Sprintf("$timestamp: %v $at: %s $message: %v", time.Now().Unix(), domain, err))
}
