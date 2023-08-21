package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	cs "github.com/flaneur4dev/good-microservices/internal/contracts"
	"github.com/redis/go-redis/v9"
)

type redStore struct {
	client *redis.Client
}

func New(a, p string) (*redStore, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     a,
		Password: p,
		DB:       0,
	})

	sc := rc.Ping(context.TODO())
	if err := sc.Err(); err != nil {
		return nil, err
	}

	fmt.Println("redis connected")
	rs := &redStore{rc}

	return rs, nil
}

func (rs *redStore) Tasks(ctx context.Context) ([]cs.Task, error) {
	zRange := rs.client.ZRange(ctx, "tasks", 0, -1)
	if err := zRange.Err(); err != nil {
		return nil, err
	}

	res, err := zRange.Result()
	if err != nil {
		return nil, err
	}

	ts := []cs.Task{}
	for _, id := range res {
		if t, err := rs.Task(ctx, id); err != nil {
			return nil, err
		} else {
			ts = append(ts, t)
		}
	}

	return ts, nil
}

func (rs *redStore) Task(ctx context.Context, id string) (cs.Task, error) {
	t := cs.Task{}

	hAll := rs.client.HGetAll(ctx, fmt.Sprintf("task:%s", id))
	if err := hAll.Err(); err != nil {
		return t, err
	}

	res, err := hAll.Result()
	if err != nil {
		return t, err
	}

	if len(res) == 0 {
		return t, errors.New("not found")
	}

	time, _ := strconv.ParseInt(res["Timestamp"], 10, 64)
	t.Id = res["Id"]
	t.Name = res["Name"]
	t.Description = res["Description"]
	t.Timestamp = time

	return t, nil
}

func (rs *redStore) Add(ctx context.Context, nt cs.Task) error {
	hSet := rs.client.HSet(ctx, fmt.Sprintf("task:%s", nt.Id), "Id", nt.Id, "Name", nt.Name, "Description", nt.Description, "Timestamp", nt.Timestamp)
	if err := hSet.Err(); err != nil {
		return err
	}

	z := redis.Z{Score: float64(nt.Timestamp), Member: nt.Id}
	zAdd := rs.client.ZAdd(ctx, "tasks", z)
	if err := zAdd.Err(); err != nil {
		return err
	}

	return nil
}

func (rs *redStore) Remove(ctx context.Context, id string) error {
	if err := rs.client.Unlink(ctx, fmt.Sprintf("task:%s", id)).Err(); err != nil {
		return err
	}

	if err := rs.client.ZRem(ctx, "tasks", id).Err(); err != nil {
		return err
	}

	return nil
}

func (rs *redStore) Check() error {
	return nil
}

// func (rs *redStore) Close() error {
// 	return nil
// }
