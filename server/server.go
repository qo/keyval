package server

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/qo/keyval/store"
)

var s *store.Store

func init() {
	st, err := store.CreateStore()
	if err != nil {
		log.Fatal("unexpected error occured during init")
	}
	s = st
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hey, this is a simple key-value store :)\n"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	val, err := s.Get(key)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := []byte("got: " + key + " - " + val + "\n")
	w.Write(res)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	valBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	val := string(valBytes)

	err = s.Put(key, val)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := []byte("put: " + key + " - " + val + "\n")
	w.Write(res)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	err := s.Delete(key)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := []byte("deleted: " + key + "\n")
	w.Write(res)
}

func Serve() {
	r := chi.NewRouter()
	r.Get("/", indexHandler)
	api := chi.NewRouter()
	api.Route("/key/{key}", func(r chi.Router) {
		r.Get("/", getHandler)
		r.Delete("/", deleteHandler)
		r.Put("/", putHandler)
	})
	r.Mount("/v1", api)
	log.Fatal(http.ListenAndServe(":8090", r))
}
