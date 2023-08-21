package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	cs "github.com/flaneur4dev/good-microservices/internal/contracts"
	hs "github.com/flaneur4dev/good-microservices/internal/handlers"
	"github.com/flaneur4dev/good-microservices/internal/storage/memory"
	"github.com/flaneur4dev/good-microservices/internal/storage/redis"
)

type Store interface {
	Tasks(ctx context.Context) ([]cs.Task, error)
	Task(ctx context.Context, id string) (cs.Task, error)
	Add(ctx context.Context, nt cs.Task) error
	Remove(ctx context.Context, id string) error
	Check() error
	// Close() error
}

const isMem = false

func main() {
	var s Store

	if isMem {
		s = memory.New()
	} else {
		r, err := redis.New(
			stringEnv("REDIS_HOST", "localhost:6379"),
			stringEnv("REDIS_PASSWORD", ""),
		)
		if err != nil {
			log.Fatal("can't connect to redis: ", err)
		}
		s = r
	}
	// defer s.Close()

	r := chi.NewRouter()

	r.Get("/ping", hs.HandleCheck(s))
	r.Get("/task", hs.HandleTasks(s))
	r.Get("/task/{id}", hs.HandleTask(s))
	r.Post("/task", hs.HandleAdd(s))
	r.Delete("/task/{id}", hs.HandleRemove(s))

	err := http.ListenAndServe(stringEnv("TASK_MANAGER_HOST", "localhost:8080"), r)
	if err != nil {
		log.Fatal(err)
	}
}

func stringEnv(name, defaultVal string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	} else {
		return defaultVal
	}
}
