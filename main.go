package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/hello", hello)
	r.Get("/panic", occurPanic)
	r.Get("/panic-recover", panicRecover)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func occurPanic(w http.ResponseWriter, r *http.Request) {
	eg := errgroup.Group{}
	eg.Go(func() error {
		panic("panic occurred")
	})
	if err := eg.Wait(); err != nil {
		log.Println(err)
	}
}

func panicRecover(w http.ResponseWriter, r *http.Request) {
	eg := errgroup.Group{}
	eg.Go(recoverer(func() {
		panic("panic occurred")
	}))
	if err := eg.Wait(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func recoverer(f func()) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v", r)
			}
		}()
		f()
		return
	}

}
