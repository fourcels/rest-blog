package models

import "time"

type User struct {
	ID        uint      `json:"id,omitempty"`
	Username  string    `gorm:"uniqueIndex" json:"username,omitempty"`
	Password  string    `json:"-"`
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
