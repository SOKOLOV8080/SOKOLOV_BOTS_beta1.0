package models

import "time"

type Product struct {
	ID          int
	SellerID    int
	Title       string
	Description string
	Price       float64
	GithubURL   string
	ContactInfo string
	FilePath    string
	CreatedAt   time.Time
}
