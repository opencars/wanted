package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Store interface {
	Revision() RevisionRepository
	Vehicle() VehicleRepository
}
