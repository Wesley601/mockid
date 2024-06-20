package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/wesley601/mockid/db"
	"github.com/wesley601/mockid/entities"
)

type RequestMatcherCached struct {
	db   *sql.DB
	data []entities.Mappings
	dao  RequestDAO
}

func NewRequestMatcherCached(db *sql.DB, dao RequestDAO) (*RequestMatcherCached, error) {
	matcher := RequestMatcherCached{
		db:  db,
		dao: dao,
	}
	paths, err := GetMappingPath()
	if err != nil {
		return nil, err
	}
	for _, v := range paths {
		dat, err := os.ReadFile(v)
		if err != nil {
			log.Printf("Unable to read the file %s: %s\n", v, err.Error())
			continue
		}
		var data entities.Mappings
		err = json.Unmarshal(dat, &data)
		if err != nil {
			return nil, err
		}
		matcher.data = append(matcher.data, data)
	}

	return &matcher, nil
}

func (m *RequestMatcherCached) Match(r *http.Request) (*Resp, error) {
	for _, data := range m.data {
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
