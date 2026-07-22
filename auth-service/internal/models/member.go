package models

import "time"

// Member represents a library member who can borrow books.
type Member struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:member" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
