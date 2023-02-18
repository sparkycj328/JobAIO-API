package data

import "time"

type Company struct {
	ID        int64     // Unique integer id for the company
	Country   string    // Company name
	Total     int       // total amount of job available
	URL       string    // URL location where resource is located
	CreatedAt time.Time // created timestamp for the data
}
