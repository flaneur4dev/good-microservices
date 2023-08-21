package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	cs "github.com/flaneur4dev/good-microservices/internal/contracts"
)

type (
	Tasks interface {
		Tasks(ctx context.Context) ([]cs.Task, error)
	}
	Task interface {
		Task(ctx context.Context, id string) (cs.Task, error)
	}
	Adder interface {
		Add(ctx context.Context, nt cs.Task) error
	}
	Remover interface {
		Remove(ctx context.Context, id string) error
	}
	Checker interface {
		Check() error
	}
)

func HandleTasks(rep Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := rep.Tasks(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v := struct {
			Tasks []cs.Task `json:"tasks"`
		}{ts}

		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleTask(rep Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		t, err := rep.Task(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		v := struct {
			Task cs.Task `json:"task"`
		}{t}

		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleAdd(rep Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var res cs.Task
		err = json.Unmarshal(reqBody, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = rep.Add(r.Context(), res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v := struct {
			Task    cs.Task `json:"task"`
			Message string  `json:"message"`
		}{res, "task created successfully"}

		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleRemove(rep Remover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := rep.Remove(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v := struct {
			Id      string `json:"id"`
			Message string `json:"message"`
		}{id, "deleted"}

		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleCheck(rep Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := rep.Check(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
