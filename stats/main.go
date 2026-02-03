package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var _ Repo = &InMemoryRepo{}

type InMemoryRepo struct {
	data map[string]int
	rw   sync.RWMutex
}

type Repo interface {
	Increment(key string)
	GetAll() map[string]int
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		data: make(map[string]int),
	}
}

func (r *InMemoryRepo) Increment(key string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.data[key]++
}

func (r *InMemoryRepo) GetAll() map[string]int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	c := make(map[string]int)
	for k, v := range r.data {
		c[k] = v
	}
	return c
}

func main() {
	mux := http.NewServeMux()
	repo := NewInMemoryRepo()
	handler := ReqMiddledleware(VisitHandler(repo))
	mux.Handle("/visit", handler)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	server := http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	go func() {
		log.Println("starting server on :8090")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server : %v", err)
		}
	}()

	<-sig
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server : %v", err)
	}
}

func ReqMiddledleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("received request: %s %s", req.Method, req.URL.Path)
		next.ServeHTTP(w, req)
	})
}

func VisitHandler(repo Repo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		repo.Increment("visits")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
