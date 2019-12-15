package worker

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/opencars/wanted/pkg/logger"
	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store"
)

var (
	ErrEmptyArr = errors.New("revision is empty")
)

type Worker struct {
	tree *redblacktree.Tree
}

type Node struct {
	Status model.Status
}

func New() *Worker {
	return &Worker{
		tree: redblacktree.NewWithStringComparator(),
	}
}

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

func (w *Worker) Parse(revision string, input io.Reader) ([]model.Vehicle, []string, error) {
	buf := bufio.NewReader(input)
	b, err := buf.Peek(16)
	if err == io.EOF {
		return nil, nil, ErrEmptyArr
	}

	if err != nil {
		return nil, nil, err
	}

	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		if _, err := buf.Discard(16); err != nil {
			return nil, nil, err
		}
	}

	dec := json.NewDecoder(buf)
	// Read the array open bracket.
	if _, err := dec.Token(); err != nil {
		return nil, nil, err
	}
	checked := make(map[string]bool)
	removedNodes := make([]string, 0)
	newTransport := make([]model.Vehicle, 0)

	// New vehicles.
	for dec.More() {
		var tmp model.WantedVehicle
		err := dec.Decode(&tmp)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, err
		}

		_, ok := w.tree.Get(tmp.ID)
		if !ok {
			v, err := model.VehicleFromGov(revision, &tmp)
			if err != nil {
				return nil, nil, err
			}
			newTransport = append(newTransport, *v)
			continue
		}

		checked[tmp.ID] = true
	}

	// Removed vehicles.
	it := w.tree.Iterator()
	for it.Next() {
		id, node := it.Key().(string), it.Value().(*Node)
		if _, ok := checked[id]; !ok && node.Status == model.StatusStolen {
			node.Status = model.StatusRemoved
			removedNodes = append(removedNodes, id)
		}
	}

	// Append new stolen vehicles to state.
	for _, v := range newTransport {
		node := &Node{Status: v.Status}
		w.tree.Put(v.ID, node)
	}

	return newTransport, removedNodes, nil
}

func (w *Worker) Load(s store.Store) error {
	vehicles, err := s.Vehicle().All()
	if err != nil {
		return err
	}

	for _, v := range vehicles {
		node := &Node{Status: v.Status}
		w.tree.Put(v.ID, node)
	}

	logger.Info("Loaded %d vehicles", len(vehicles))

	return nil
}
