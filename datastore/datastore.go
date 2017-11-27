package datastore

// Provide default options with sensible values
var DefaultListOptions = ListOptions{
	PerPage: 20,
	Page:    0,
}

// Parameters for list, can be serialized to and from a
// URL or JSON.
type ListOptions struct {
	PerPage int `schema:"per_page"`
	Page    int `schema:"page"`
}
