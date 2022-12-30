package worker

import (
	"encoding/json"
	"io"

	"github.com/emirpasic/gods/trees/redblacktree"

	"github.com/opencars/seedwork/logger"
	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/domain"
	"github.com/opencars/wanted/pkg/domain/model"
)

type Worker struct {
	tree *redblacktree.Tree
}

func New() *Worker {
	return &Worker{
		tree: redblacktree.NewWithStringComparator(),
	}
}

func (w *Worker) Parse(revision string, input io.Reader) ([]model.Vehicle, []string, error) {
	input, err := bom.NewReader(input)
	if err != nil {
		return nil, nil, err
	}

	dec := json.NewDecoder(input)
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

		tmp.CheckSum = tmp.CalculateCheckSum()
		_, ok := w.tree.Get(tmp.CheckSum)
		if !ok {
			v, err := model.VehicleFromGov(revision, &tmp)
			if err != nil {
				return nil, nil, err
			}
			newTransport = append(newTransport, *v)
			continue
		}

		checked[tmp.CheckSum] = true
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
		w.tree.Put(v.CheckSum, NewNode(v.Status))
	}

	return newTransport, removedNodes, nil
}

func (w *Worker) Load(s domain.Store) error {
	vehicles, err := s.Vehicle().All()
	if err != nil {
		return err
	}

	for _, v := range vehicles {
		w.tree.Put(v.CheckSum, NewNode(v.Status))
	}

	logger.Infof("Loaded %d vehicles", len(vehicles))

	return nil
}
