package data

import (
	"github.com/sparkycj328/JobAIO-API/internal/validator"
	"time"
)

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

// ValidateCompany will perform validation checks on each
func ValidateCompany(v *validator.Validator, company *Company) {
	v.Check(company.Name != "", "name", "must be provided")
	v.Check(len(company.Name) <= 100, "name", "must not be more than 100 bytes long")

	for _, country := range company.Countries {
		v.Check(country.Country != "", "country", "must be provided")
		v.Check(len(country.Country) <= 100, "country", "must not be more than 100 bytes long")

		v.Check(country.Total < 0, "amount", "cannot be a negative number")

		v.Check(country.URL != "", "url", "must be provided")
		v.Check(len(country.URL) <= 100, "url", "must not be more than 200 bytes long")

		v.Check(country.CreatedAt.Before(time.Now()) || country.CreatedAt.After(time.Now()), "created",
			"must be equal to today's date")

	}
}
