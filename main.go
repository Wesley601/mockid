package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/local.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer db.Close()
	if err != nil {
		panic(err)
	}
	http.Handle("/", &MapHandler{matcher: &RequestMatcherLive{db: db}})

	http.ListenAndServe(":3000", nil)
}

var errRespNotFound = errors.New("response not found")

type RequestMatcher interface {
	Match(url, method string) (*Resp, error)
}

type MapHandler struct {
	matcher RequestMatcher
}

func (m *MapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	method := r.Method
	resp, err := m.matcher.Match(url, method)
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
