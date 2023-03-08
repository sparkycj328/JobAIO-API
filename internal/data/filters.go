package data

import "github.com/sparkycj328/JobAIO-API/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// ValidateFilters will perform validation checks on the filter query paramters
// provided by the http client in their url
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "page must be greater than 0")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum 0f 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Iterate through the sortsafelist and validated that the sort parameter
	// is a valid sort list
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
