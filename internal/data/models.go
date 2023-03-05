package data

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound will be used by our Get function if an error occurs
// while retrieving a set of records
var (
	ErrRecordNotFound = errors.New("records not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models wraps the VendorModel struct and will wrap other necessary structs in the future
type Models struct {
	Vendors VendorModel
}

// NewModel returns a Models struct containing the initialized VendorModel
func NewModel(db *sql.DB) Models {
	return Models{
		Vendors: VendorModel{DB: db},
	}
}
