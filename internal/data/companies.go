package data

import (
	"database/sql"
	"errors"
	"github.com/sparkycj328/JobAIO-API/internal/validator"
	"time"
)

type Company struct {
	Name      string     `json:"company,omitempty"` // company name
	ID        int64      `json:"-"`                 // Unique integer id for the company
	Country   string     `json:"country"`           // Country name
	Total     int        `json:"total"`             // total amount of job available
	URL       string     `json:"url"`               // URL location where resource is located
	CreatedAt *time.Time `json:"created,omitempty"` // created timestamp for the data
}

// ValidateCompany will perform validation checks on each
func ValidateCompany(v *validator.Validator, c *Company) {
	v.Check(c.Name != "", "name", "must be provided")
	v.Check(len(c.Name) <= 100, "name", "must not be more than 100 bytes long")
	v.Check(c.Country != "", "country", "must be provided")
	v.Check(len(c.Country) <= 100, "country", "must not be more than 100 bytes long")
	v.Check(c.Total >= 0, "amount", "cannot be a negative number")
	v.Check(c.URL != "", "url", "must be provided")
	v.Check(len(c.URL) <= 100, "url", "must not be more than 200 bytes long")
}

// VendorModel wraps the sql.DB connection pool in a struct
type VendorModel struct {
	DB *sql.DB
}

// Insert will take the company struct and insert the data into our database
func (m *VendorModel) Insert(c *Company) error {
	query := `
			INSERT INTO jobs (vendor, country, amount, url)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at`
	args := []any{c.Name, c.Country, c.Total, c.URL}
	return m.DB.QueryRow(query, args...).Scan(&c.ID, &c.CreatedAt)
}

// GetRows will for fetching specific records from the jobs table
func (m *VendorModel) GetRows(vendor string) (*[]Company, error) {
	// if vendor string is empty return an error
	if vendor == "" {
		return nil, ErrRecordNotFound
	}
	// define a slice of company struct which will
	// be used to store the rows queried
	countries := make([]Company, 0)

	// define the SQL statement
	query := `SELECT id, created_at, country, amount, url
		  		FROM jobs
				WHERE vendor = $1 AND created_at::date = CURRENT_DATE AND amount > 0 ORDER BY country;`
	rows, err := m.DB.Query(query, vendor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		country := Company{}

		err = rows.Scan(&country.ID, &country.CreatedAt, &country.Country, &country.Total, &country.URL)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
			}
		}
		countries = append(countries, country)
	}
	return &countries, nil
}

// Update will update the specified records in the job table
func (m *VendorModel) Update(c *Company) error {

	return nil
}

func (m *VendorModel) Delete(vendor string) error {

	return nil
}
