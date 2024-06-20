package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/wesley601/mockid/db"
	"github.com/wesley601/mockid/views"
)

type RequestDAO interface {
	FindById(id string) (*db.RequestSaved, error)
	Flush() error
	List() ([]db.RequestSaved, error)
}

type RequestHandler struct {
	db  *sql.DB
	dao RequestDAO
}

func NewRequestHandler(db *sql.DB, dao RequestDAO) *RequestHandler {
	return &RequestHandler{
		db:  db,
		dao: dao,
	}
}

func (re *RequestHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	request, err := re.dao.FindById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("no request with id: %s", id), http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	views.Show(*request).Render(r.Context(), w)
}

func (re *RequestHandler) Index(w http.ResponseWriter, r *http.Request) {
	requests, err := re.dao.List()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	views.Base(views.Home(requests)).Render(r.Context(), w)
}

func (re *RequestHandler) Flush(w http.ResponseWriter, r *http.Request) {
	err := re.dao.Flush()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
