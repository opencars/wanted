package event

import "encoding/json"

type Message struct {
	Kind    Kind            `json:"kind"`
	Payload json.RawMessage `json:"payload"`
}

type Kind string

const (
	KindRevision Kind = "revision"
)

type Revision struct {
	ID        string `json:"id"`
	Resource  string `json:"resource"`
	Package   string `json:"package"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

func New(kind Kind, object interface{}) (*Message, error) {
	payload, err := json.Marshal(&object)
	if err != nil {
		return nil, err
	}

	return &Message{
		Kind:    kind,
		Payload: payload,
	}, nil
}
