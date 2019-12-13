package model

import (
	"strings"
	"time"

	"github.com/opencars/govdata"
)

// Revision represents storage model for revision entity.
type Revision struct {
	ID          string    `json:"id" db:"id"`
	URL         string    `json:"url" db:"url"`
	FileHashSum *string   `json:"file_hash_sum" db:"file_hash_sum"`
	Removed     int       `json:"removed" db:"removed"`
	Added       int       `json:"added" db:"added"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// RevisionStatMonth represents revisions aggregate by month.
type RevisionStatMonth struct {
	Month   time.Month `db:"month" json:"month"`
	Year    int        `db:"year" json:"year"`
	Added   int        `db:"added" json:"added"`
	Removed int        `db:"removed" json:"removed"`
}

func RevisionFromGov(revision *govdata.Revision) *Revision {
	parts := strings.Split(revision.URL, "/")

	return &Revision{
		ID:          parts[len(parts)-1],
		URL:         revision.URL,
		FileHashSum: revision.FileHashSum,
		CreatedAt:   revision.ResourceCreated.Time,
	}
}
