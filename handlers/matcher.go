package handlers

import (
	"errors"
	"net/http"

	"github.com/wesley601/mockid/services"
)

type RequestMatcher interface {
	Match(r *http.Request) (*services.Resp, error)
}

type MapHandler struct {
	matcher RequestMatcher
}

func NewMapHandler(m RequestMatcher) *MapHandler {
	return &MapHandler{
		matcher: m,
	}
}

func (m *MapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := m.matcher.Match(r)
	if err != nil {
		if errors.Is(err, services.ErrRespNotFound) {
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
