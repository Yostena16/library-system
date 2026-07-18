package models

import "time"

// Fine represents a penalty for returning a book late.
type Fine struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LoanID    uint      `gorm:"not null" json:"loan_id"`
	MemberID  uint      `gorm:"not null" json:"member_id"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Paid      bool      `gorm:"default:false" json:"paid"`
	CreatedAt time.Time `json:"created_at"`
}
