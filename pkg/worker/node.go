package worker

import (
	"github.com/opencars/wanted/pkg/model"
)

type Node struct {
	Status model.Status
}

func NewNode(status model.Status) *Node {
	return &Node{
		Status: status,
	}
}
