package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/local.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	m, err := NewMigratorEmbed(db)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
	http.Handle("/", &MapHandler{matcher: &RequestMatcherLive{db: db}})

	log.Println("Server starting at :3000")
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

var errRespNotFound = errors.New("response not found")

type RequestMatcher interface {
	Match(r *http.Request) (*Resp, error)
}

type MapHandler struct {
	matcher RequestMatcher
}

func (m *MapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := m.matcher.Match(r)
	if err != nil {
		if errors.Is(err, errRespNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message":"no mapping found"}`))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(int(resp.Status))
	for k, v := range resp.Headers {
		w.Header().Add(k, v)
	}
	w.Write(resp.Body)
}
