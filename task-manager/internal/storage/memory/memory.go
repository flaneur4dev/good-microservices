package memory

import (
	"context"
	"errors"

	cs "github.com/flaneur4dev/good-microservices/internal/contracts"
)

type memStore struct {
	tasks map[string]cs.Task
}

func New() *memStore {
	return &memStore{
		tasks: map[string]cs.Task{},
	}
}

func (ms *memStore) Tasks(_ context.Context) ([]cs.Task, error) {
	ts := make([]cs.Task, 0, len(ms.tasks))
	for _, v := range ms.tasks {
		ts = append(ts, v)
	}
	return ts, nil
}

func (ms *memStore) Task(_ context.Context, id string) (cs.Task, error) {
	t, ok := ms.tasks[id]
	if !ok {
		return cs.Task{}, errors.New("not found")
	}
	return t, nil
}

func (ms *memStore) Add(_ context.Context, nt cs.Task) error {
	ms.tasks[nt.Id] = nt
	return nil
}

func (ms *memStore) Remove(_ context.Context, id string) error {
	delete(ms.tasks, id)
	return nil
}

func (ms *memStore) Check() error {
	return nil
}

// func (ms *memStore) Close() error {
// 	return nil
// }
