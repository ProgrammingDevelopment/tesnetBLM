package models

import "time"

type User struct {
	ID        string    `json:"id"`
	NIK       string    `json:"nik"`
	Nama      string    `json:"nama"`
	Whatsapp  string    `json:"whatsapp"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // Encrypted
	CreatedAt time.Time `json:"created_at"`
}
