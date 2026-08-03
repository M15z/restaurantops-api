// internal/models/menu.go
package models

import "time"

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MenuItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	Category    Category  `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}
