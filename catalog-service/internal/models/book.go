package models

// Book represents a single title in the library's catalog.
type Book struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Title           string `gorm:"not null" json:"title"`
	Author          string `json:"author"`
	Category        string `json:"category"`
	ISBN            string `gorm:"unique" json:"isbn"`
	TotalCopies     int    `json:"total_copies"`
	AvailableCopies int    `json:"available_copies"`
}
