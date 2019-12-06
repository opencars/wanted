package worker

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store"
)

var (
	ErrEmptyArr = errors.New("revision is empty")
)

type Worker struct {
	state store.Transport
}

func New() *Worker {
	return &Worker{
		state: make(store.Transport, 0),
	}
}

func (w *Worker) Parse(revision string, input io.Reader) ([]model.Vehicle, int, int, error) {
	buf := bufio.NewReader(input)

	_, err := buf.Peek(32)
	if err == io.EOF {
		fmt.Printf("SKIPPED: %s\n", revision)
		return nil, 0, 0, ErrEmptyArr
	}

	if err != nil {
		return nil, 0, 0, err
	}

	dec := json.NewDecoder(input)

	// Read the array open bracket.
	if _, err := dec.Token(); err != nil {
		return nil, 0, 0, err
	}

	checked := make(map[string]bool)
	result := make([]model.Vehicle, 0)
	newTransport := make([]model.Vehicle, 0)
	added := 0
	removed := 0

	// New vehicles.
	for dec.More() {
		var tmp model.WantedVehicle
		err := dec.Decode(&tmp)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, 0, err
		}

		i := w.state.Search(tmp.ID)

		if i == -1 /* New stolen vehicle */ {
			v, err := model.VehicleFromGov(revision, &tmp)
			if err != nil {
				return nil, 0, 0, err
			}
			v.RevisionID = revision
			added++
			newTransport = append(newTransport, *v)
			continue
		}

		checked[w.state[i].ID] = true
		if w.state[i].Status != model.StatusStolen {
			added++
			w.state[i].Status = model.StatusStolen
			result = append(result, w.state[i])
		}
	}

	// Removed vehicles.
	for i, v := range w.state {
		if _, ok := checked[v.ID]; !ok && v.Status == model.StatusStolen {
			removed++
			w.state[i].Status = model.StatusRemoved
			result = append(result, w.state[i])
		}
	}

	// Append new stolen vehicles to state.
	for _, v := range newTransport {
		w.state = append(w.state, v)
		result = append(result, v)
	}

	// Sort newly updated array.
	sort.Sort(w.state)
	fmt.Println("Sorted ", sort.IsSorted(w.state))

	return result, added, removed, nil
}

func (w *Worker) Load(s store.Store) error {
	vehicles, err := s.Vehicle().All()
	if err != nil {
		return err
	}

	log.Printf("Loading %d vehicles\n", len(vehicles))
	tmp := store.Transport(vehicles)
	sort.Sort(tmp)
	w.state = tmp

	return nil
}
