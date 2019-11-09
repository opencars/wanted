package govdata

import (
	"context"
	"log"
	"time"
)

type Subscription struct {
	id      string
	current []Revision
}

func NewSubscription(id string) *Subscription {
	return &Subscription{}
}

func listen(id string, ids []string, revisions chan<- Revision) {
	for {
		resource, err := DefaultClient.ResourceShow(context.Background(), id)
		if err != nil {
			log.Println("ResourceShow:", err)
			<-time.After(30 * time.Second)
			continue
		}

		for i := len(resource.Revisions) - 1 - len(ids); i >= 0; i-- {
			revisions <- resource.Revisions[i]
			ids = append(ids, resource.Revisions[i].ID)
		}

		<-time.After(3 * time.Minute)
	}
}

func Subscribe(id string, ids ...string) <-chan Revision {
	revisions := make(chan Revision)
	go listen(id, ids, revisions)
	return revisions
}
