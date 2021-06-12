package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Ticket struct {
	gorm.Model
	Status string `json:"status" gorm:"default:'open'"`
	UserId uint   `json:"user"`
}

// ENDPOINTS

func createTicket(w http.ResponseWriter, r *http.Request) {
	claims, valid := extractClaims(r.Header.Get("Authorization"))
	var currentUser User
	db.Where("Username = ?", claims["username"]).First(&currentUser)

	if !valid {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	ticket := Ticket{
		UserId: currentUser.ID,
	}
	db.Save(&ticket)
	json.NewEncoder(w).Encode(ticket)
}

func getTickets(w http.ResponseWriter, r *http.Request) {
	var tickets []Ticket
	db.Order("ID DESC").Find(&tickets)
	json.NewEncoder(w).Encode(tickets)
}

func getTicket(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()["id"][0]
	var ticket Ticket
	db.First(&ticket, id)
	json.NewEncoder(w).Encode(ticket)
}

func deleteTicket(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()["id"][0]
	var ticket Ticket
	db.First(&ticket, id)
	db.Unscoped().Where("ID = ?", id).Delete(&Ticket{})
	json.NewEncoder(w).Encode(ticket)

}

func isValidStatus(status string) bool {
	if status != "open" && status != "closed" {
		return false
	} else {
		return true
	}
}

func updateTicket(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()["id"][0]
	var updatedTicket Ticket
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &updatedTicket)

	if isValidStatus(updatedTicket.Status) {
		db.First(&updatedTicket, id)
		db.Model(&updatedTicket).Updates(map[string]interface{}{"status": updatedTicket.Status})
		json.NewEncoder(w).Encode(updatedTicket)
	} else {
		http.Error(w, "Invalid session", http.StatusUnprocessableEntity)

	}

}
