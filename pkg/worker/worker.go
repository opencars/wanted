package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"

	"github.com/opencars/wanted/pkg/storage"
)

type Worker struct {
	state storage.Transport
}

func New() *Worker {
	return &Worker{make(storage.Transport, 0)}
}

func (w *Worker) Parse(revision string, input io.Reader) ([]storage.WantedVehicle, int, int, error) {
	dec := json.NewDecoder(input)

	// Read the array open bracket.
	if _, err := dec.Token(); err != nil {
		return nil, 0, 0, err
	}

	checked := make(map[string]bool)
	result := make([]storage.WantedVehicle, 0)
	newTransport := make([]storage.WantedVehicle, 0)
	added := 0
	removed := 0

	// New vehicles.
	for dec.More() {
		var tmp storage.WantedVehicle
		err := dec.Decode(&tmp)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, 0, err
		}

		i := w.state.Search(tmp.ID)

		if i == -1 /* New stolen vehicle */ {
			tmp.Status = storage.StatusStolen
			tmp.RevisionID = revision
			added++
			newTransport = append(newTransport, tmp)
			continue
		}

		checked[w.state[i].ID] = true
		if w.state[i].Status != storage.StatusStolen {
			added++
			w.state[i].Status = storage.StatusStolen
			result = append(result, w.state[i])
		}
	}

	// Removed vehicles.
	for i, v := range w.state {
		if _, ok := checked[v.ID]; !ok && v.Status == storage.StatusStolen {
			removed++
			w.state[i].Status = storage.StatusRemoved
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

	return result, added, removed, nil
}

func (w *Worker) Load(store *storage.Store) error {
	vehicles, err := store.AllVehicles()
	if err != nil {
		return err
	}

	fmt.Println(len(vehicles))
	tmp := storage.Transport(vehicles)
	sort.Sort(tmp)
	w.state = tmp

	return nil
}

func (w *Worker) Fix(vehicles []storage.WantedVehicle) {
	for i := range vehicles {
		vehicles[i].Brand, vehicles[i].Kind = storage.ParseKind(vehicles[i].Brand)

		// Remove unnecessary lexemes from vehicle kind.
		for _, lexeme := range []string{"АВТОБУС ", "АВТОТРАНСПОРТ"} {
			vehicles[i].Kind = strings.ReplaceAll(strings.ToUpper(vehicles[i].Kind), lexeme, "")
		}

		vehicles[i].Kind = strings.TrimFunc(vehicles[i].Kind, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		// Removed redundant spaces.
		vehicles[i].OVD = strings.TrimSpace(vehicles[i].OVD)
		vehicles[i].Color = strings.ToUpper(strings.TrimSpace(vehicles[i].Color))
		vehicles[i].Brand = strings.ToUpper(strings.TrimSpace(vehicles[i].Brand))

		// TODO: Transliterate number into cyrillic.
		// vehicles[i].Number =

		vehicles[i].BodyNumber = strings.TrimSpace(vehicles[i].BodyNumber)
		vehicles[i].ChassisNumber = strings.TrimSpace(vehicles[i].ChassisNumber)
		vehicles[i].EngineNumber = strings.TrimSpace(vehicles[i].EngineNumber)
	}
}
