package db

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
}
