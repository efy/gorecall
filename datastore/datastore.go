package datastore

import "fmt"

var (
	ErrInvalidDB = fmt.Errorf("invalid database")
)

// Validator specifies the method Validate
// for perfoming validation and returning errors
type Validator interface {
	Validate() (bool, []error)
}

// Provide default options with sensible values
var DefaultListOptions = ListOptions{
	PerPage: 20,
	Page:    0,
	OrderBy: "created",
	Order:   "DESC",
}

// Parameters for list, can be serialized to and from a
// URL or JSON.
type ListOptions struct {
	PerPage int    `schema:"per_page"`
	Page    int    `schema:"page"`
	OrderBy string `schema:"order_by"`
	Order   string `schema:"order"`
}
