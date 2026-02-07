package models

import "time"

type Ticket struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	LocationID   string    `json:"location_id"`
	LocationName string    `json:"location_name"`
	TicketNumber string    `json:"ticket_number"`
	Code         string    `json:"code"`
	TimeSlot     string    `json:"time_slot"`
	CreatedAt    time.Time `json:"created_at"`
}

type Location struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Quota int64  `json:"quota"`
	Region string `json:"region"`
}
