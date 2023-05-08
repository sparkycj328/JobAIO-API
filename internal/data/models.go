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
	Users   UserModel
}

// NewModel returns a Models struct containing the initialized VendorModel and UsersModel
func NewModel(db *sql.DB) Models {
	return Models{
		Vendors: VendorModel{DB: db},
		Users:   UserModel{DB: db},
	}
}
