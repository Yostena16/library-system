package models

import "time"

// Loan represents a book borrowed by a member.
type Loan struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	MemberID   uint       `gorm:"not null" json:"member_id"`
	BookID     uint       `gorm:"not null" json:"book_id"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueDate    time.Time  `json:"due_date"`
	ReturnedAt *time.Time `json:"returned_at"`
	Status     string     `gorm:"default:borrowed" json:"status"`
}
