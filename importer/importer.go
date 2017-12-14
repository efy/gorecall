// Package importer handles persisting items to the database and generating a status report
// of successfully imported items returning a list of failed imports and their respective errors
package importer

import "github.com/efy/gorecall/datastore"

type ItemReport struct {
	Bookmark datastore.Bookmark `json:"bookmark"`
	Error    error              `json:"error"`
}

type Report struct {
	SuccessCount int          `json:"success_count"`
	FailureCount int          `json:"failure_count"`
	Results      []ItemReport `json:"results"`
}

type Options struct {
	RunWebinfo bool
}

var DefaultOptions = Options{
	RunWebinfo: false,
}

func Import(src []datastore.Bookmark, dst datastore.BookmarkRepo, opts Options) (*Report, error) {
	report := &Report{}

	for _, bm := range src {
		item := ItemReport{}
		b, err := dst.Create(&bm)
		if err != nil {
			b = &bm
			item.Error = err
			report.FailureCount++
		} else {
			report.SuccessCount++
		}
		item.Bookmark = *b

		report.Results = append(report.Results, item)
	}

	return report, nil
}
