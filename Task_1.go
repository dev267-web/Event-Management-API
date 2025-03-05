# Event-Management-API
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Event struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	Organiser string `json:"organiser"`
	Status    string `json:"status"`
}
type Ticket struct {
	ID      int    `json:"id"`
	EventID int    `json:"event_id"`
	User    string `json:"user"`
	SeatNO  string `json:"seat_no"`
	Status  string `json:"status"`
}

var events []Event
var tickets []Ticket

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent Event

	json.NewDecoder(r.Body).Decode(&newEvent)
	newEvent.ID = len(events) + 1
	newEvent.Status = "Pending"
	events = append(events, newEvent)
	w.Header().Set("Content-Type",
		"application/json")
	json.NewEncoder(w).Encode(newEvent)
}
func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
func getEventByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i := range events {
		if events[i].ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(events)
			return
		}
	}
	http.Error(w, "Event not found",
		http.StatusNotFound)
}
func approveEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i := range events {
		if events[i].ID == id {
			events[i].Status = "Approved"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(events)
			return
		}
	}
	http.Error(w, "Event not found",
		http.StatusNotFound)
}
func rejectEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i := range events {
		if events[i].ID == id {
			events[i].Status = "Rejected"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(events)
			return
		}
	}
	http.Error(w, "Event not found",
		http.StatusNotFound)
}

func generateTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i := range events {
		if events[i].ID == id {
			if events[i].Status != "Approved" {
				http.Error(w, "Event not Approved yet", http.StatusForbidden)
				return
			}
			var newTicket Ticket

			json.NewDecoder(r.Body).Decode(&newTicket)
			newTicket.ID = len(tickets) + 1
			newTicket.EventID = id
			newTicket.SeatNO = "A" +
				strconv.Itoa(rand.Intn(100))
			tickets = append(tickets, newTicket)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newTicket)
			return
		}
		http.Error(w, "Event not found", http.StatusNotFound)

	}
}
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/events", createEvent).Methods("POST")
	r.HandleFunc("/events", getEvents).Methods("GET")
	r.HandleFunc("/events/{id}", getEventByID).Methods("GET")
	r.HandleFunc("/events/{id}/approve", approveEvent).Methods("PUT")
	r.HandleFunc("/events/{id}/reject", rejectEvent).Methods("PUT")
	r.HandleFunc("/events/{id}/tickets", generateTicket).Methods("POST")

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)

}
