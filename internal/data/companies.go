package data

import "time"

type Company struct {
	Name      string    `json:"company"`   // company name
	Countries []Country `json:"countries"` // job postings
}

type Country struct {
	ID        int64     `json:"id"`      // Unique integer id for the company
	Country   string    `json:"country"` // Country name
	Total     int       `json:"total"`   // total amount of job available
	URL       string    `json:"url"`     // URL location where resource is located
	CreatedAt time.Time `json"-"`        // created timestamp for the data
}
