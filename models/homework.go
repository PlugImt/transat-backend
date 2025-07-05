package models

import "time"

type Homework struct {
	ID          int       `json:"id"`
	Author      string    `json:"author"`
	CourseName  string    `json:"course_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    string    `json:"deadline"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
