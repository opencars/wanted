package store

import (
	"github.com/opencars/wanted/pkg/model"
)

// Transport is a wrapper for slice of WantedVehicle.
type Transport []model.Vehicle

// Len is the number of elements in the collection.
func (t Transport) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t Transport) Less(i, j int) bool {
	return t[i].CheckSum < t[j].CheckSum
}

// Swap swaps the elements with indexes i and j.
func (t Transport) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Search returns id of element, which has specified ID.
// Transport array should be pre-sorted.
//
// Uses binary-search algorithm under the hood.
// More about algorithm: https://en.wikipedia.org/wiki/Binary_search_algorithm.
func (t Transport) Search(id string) int {
	return t.search(0, len(t), id)
}

func (t Transport) search(from, to int, id string) int {
	pivot := (to-from)/2 + from

	if to-from <= 0 {
		return -1
	}

	if id == t[pivot].CheckSum {
		return pivot
	}

	if id > t[pivot].CheckSum {
		return t.search(pivot+1, to, id)
	}

	return t.search(from, pivot, id)
}
