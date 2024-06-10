package main

import (
	"errors"
	"net/http"
)

func main() {
	matcher, err := NewResquetMatcherLive()
	if err != nil {
		panic(err)
	}
	http.Handle("/", &MapHandler{matcher: matcher})

	http.ListenAndServe(":8080", nil)
}

var errRespNotFound = errors.New("response not found")

type Resp struct {
	Headers map[string]string
	Body    []byte
	Status  int
}

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
