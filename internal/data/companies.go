package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sparkycj328/JobAIO-API/internal/validator"
	"time"
)

type Company struct {
	ID        int64      `json:"-"`                 // Unique integer id for the company
	Name      string     `json:"company,omitempty"` // company name
	Country   string     `json:"country"`           // Country name
	Total     int        `json:"total"`             // total amount of job available
	URL       string     `json:"url"`               // URL location where resource is located
	Version   int32      `json:"version"`           // updated each time a record is updated
	CreatedAt *time.Time `json:"created,omitempty"` // created timestamp for the data
}

// ValidateCompany will perform validation checks on each field of the given Company struct
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
// acts as our POST endpoint
func (m *VendorModel) Insert(c *Company) error {
	query := `
			INSERT INTO jobs (vendor, country, amount, url)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`
	args := []any{c.Name, c.Country, c.Total, c.URL}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&c.ID, &c.CreatedAt, &c.Version)
}

// GetRecord queries our jobs table for an individual row
// this row is called using the id parameter from the URL request
func (m *VendorModel) GetRecord(id int64) (*Company, error) {

	// one last validation check
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// build the single query
	query := `
			SELECT id, created_at, vendor, country, amount, url
			FROM jobs
			WHERE id = $1`

	var record Company

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query for the matching id and based on type of error
	// return our ErrRecordNotFound error or return other error
	if err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.CreatedAt,
		&record.Name,
		&record.Country,
		&record.Total,
		&record.URL,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &record, nil
}

// GetAllRows will be used to grab all rows from the jobs table
func (m *VendorModel) GetAllRows(vendor, country string, total int, filters Filters) (*[]Company, error) {
	// define a slice of company struct which will
	// be used to store the rows queried
	countries := make([]Company, 0)

	// define the SQL statement
	query := `SELECT id, created_at, country, amount, url, version
		  	  FROM jobs
		      ORDER BY id`

	//
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate through each row returned by our query to the jobs table
	for rows.Next() {
		// declare a local instance of our company struct
		country := Company{}

		// scan the individual record values for the current row into our local struct
		// based on type of error, return different error messages
		if err := rows.Scan(
			&country.ID,
			&country.CreatedAt,
			&country.Country,
			&country.Total,
			&country.URL,
			&country.Version,
		); err != nil {
			return nil, err
		}
		// append the filled struct to our slice of rows queried.
		countries = append(countries, country)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &countries, nil
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
	query := `SELECT id, created_at, country, amount, url, version
		  		FROM jobs
				WHERE vendor = $1 AND created_at::date = CURRENT_DATE AND amount > 0 ORDER BY country`

	//
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, vendor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate through each row returned by our query to the jobs table
	for rows.Next() {
		// declare a local instance of our company struct
		country := Company{}

		// scan the individual record values for the current row into our local struct
		// based on type of error, return different error messages
		if err := rows.Scan(
			&country.ID,
			&country.CreatedAt,
			&country.Country,
			&country.Total,
			&country.URL,
			&country.Version,
		); err != nil {
			return nil, err
		}
		// append the filled struct to our slice of rows queried.
		countries = append(countries, country)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(countries) == 0 {
		return nil, ErrRecordNotFound
	}
	return &countries, nil
}

// Update will update the specified records in the job table
func (m *VendorModel) Update(c *Company) error {
	// create the prepared statement
	query := `
			UPDATE jobs
			SET vendor = $1, country = $2, amount= $3, url= $4, version = version + 1 
			WHERE id = $5 and version = $6
			RETURNING version
`
	args := []any{
		c.Name,
		c.Country,
		c.Total,
		c.URL,
		c.ID,
		c.Version,
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute the query in our jobs table
	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&c.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete will delete a record from our jobs table
// if the matching record exists
func (m *VendorModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the prepared statement
	query := `
			DELETE FROM jobs
			WHERE id  = $1
`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// we are using DB.EXEC due to not wanting any rows returned
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
