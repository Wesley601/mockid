package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/wesley601/mockid/db"
	"github.com/wesley601/mockid/entities"
)

var ErrRespNotFound = errors.New("response not found")

type Resp struct {
	Headers map[string]string
	Body    []byte
	Status  int
}

type RequestDAO interface {
	Create(db.RequestSaved) error
}

type RequestMatcherLive struct {
	db  *sql.DB
	dao RequestDAO
}

func NewRequestMatcherLive(db *sql.DB, dao RequestDAO) *RequestMatcherLive {
	return &RequestMatcherLive{
		db:  db,
		dao: dao,
	}
}

func (m *RequestMatcherLive) Match(r *http.Request) (*Resp, error) {
	paths, err := GetMappingPath()
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		file, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Unable to read the file %s: %s\n", path, err.Error())
			continue
		}

		var data entities.Mappings
		err = json.Unmarshal(file, &data)
		if err != nil {
			return nil, err
		}

		for _, mapping := range data.Mappings {
			if mapping.Request.Match(r) {
				response, err := mapping.Response.GetBody()
				if err != nil {
					return nil, err
				}

				m.dao.Create(db.RequestSaved{
					RequestedPath:   r.URL.Path,
					RequestedMethod: r.Method,
					MatchedPath:     mapping.Request.URLPattern,
					ResponseBody:    string(response),
					ResponseStatus:  mapping.Response.Status,
				})
				return &Resp{
					Status:  mapping.Response.Status,
					Headers: mapping.Response.Headers,
					Body:    response,
				}, nil
			}
		}
	}

	return nil, ErrRespNotFound
}
