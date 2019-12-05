package worker

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"unicode"

	"github.com/opencars/translit"
	"github.com/opencars/wanted/pkg/storage"
)

var (
	ErrEmptyArr = errors.New("revision is empty")
)

type Worker struct {
	state storage.Transport
}

func New() *Worker {
	return &Worker{make(storage.Transport, 0)}
}

func (w *Worker) Parse(revision string, input io.Reader) ([]storage.WantedVehicle, int, int, error) {
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

	log.Printf("Loading %d vehicles\n", len(vehicles))
	tmp := storage.Transport(vehicles)
	sort.Sort(tmp)
	w.state = tmp

	return nil
}

func TrimNilStr(lexeme *string) *string {
	if lexeme == nil {
		return nil
	}

	str := *lexeme
	str = strings.TrimSpace(str)

	str = strings.TrimFunc(str, func(r rune) bool {
		return r == '-' || r == '%' || r == '*' || r == '.'
	})

	if str == "" {
		return nil
	}

	return &str
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
		vehicles[i].Brand = strings.ToUpper(strings.TrimSpace(vehicles[i].Brand))

		// Fix Number if not nil.
		if vehicles[i].Color != nil {
			vehicles[i].Color = TrimNilStr(vehicles[i].Color)
			*vehicles[i].Color = strings.ReplaceAll(strings.ToUpper(*vehicles[i].Color), "НЕВИЗНАЧЕНИЙ", "")
			*vehicles[i].Color = strings.ReplaceAll(strings.ToUpper(*vehicles[i].Color), "НЕОПРЕДЕЛЕН", "")

			if *vehicles[i].Color == "" {
				vehicles[i].Color = nil
			}
		}

		// Transliterate number into cyrillic.
		vehicles[i].Number = TrimNilStr(vehicles[i].Number)

		// Fix color if number nil.
		if vehicles[i].Number != nil {
			*vehicles[i].Number = translit.ToUA(*vehicles[i].Number)

			if *vehicles[i].Number == "" {
				vehicles[i].Number = nil
			}
		}

		vehicles[i].BodyNumber = TrimNilStr(vehicles[i].BodyNumber)
		vehicles[i].ChassisNumber = TrimNilStr(vehicles[i].ChassisNumber)
		vehicles[i].EngineNumber = TrimNilStr(vehicles[i].EngineNumber)

		// Fix theft_date.
		vehicles[i].TheftDate = vehicles[i].TheftDate[0:10]
	}
}
