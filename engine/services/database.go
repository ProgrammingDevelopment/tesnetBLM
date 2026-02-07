package services

import (
	"encoding/json"
	"os"
	"sync"
	"war-ticket-engine/models"
)

type Database struct {
	Users   map[string]models.User   `json:"users"`
	Tickets map[string]models.Ticket `json:"tickets"`
	mu      sync.RWMutex
	path    string
}

var DB *Database

func InitDatabase(path string) *Database {
	DB = &Database{
		Users:   make(map[string]models.User),
		Tickets: make(map[string]models.Ticket),
		path:    path,
	}
	DB.Load()
	return DB
}

func (db *Database) Load() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return nil
	}

	return json.Unmarshal(data, db)
}

func (db *Database) Save() error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}

func (db *Database) GetUser(id string) (models.User, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	user, ok := db.Users[id]
	return user, ok
}

func (db *Database) SetUser(user models.User) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Users[user.ID] = user
	go db.Save()
}

func (db *Database) GetUserByEmailOrPhone(identifier string) (models.User, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, u := range db.Users {
		if u.Email == identifier || u.Whatsapp == identifier {
			return u, true
		}
	}
	return models.User{}, false
}

func (db *Database) GetUserByNIK(nik string) (models.User, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, u := range db.Users {
		if u.NIK == nik {
			return u, true
		}
	}
	return models.User{}, false
}

func (db *Database) GetTicket(id string) (models.Ticket, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	ticket, ok := db.Tickets[id]
	return ticket, ok
}

func (db *Database) SetTicket(ticket models.Ticket) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Tickets[ticket.ID] = ticket
	go db.Save()
}
