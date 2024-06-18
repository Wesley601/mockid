package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	var files [][]byte
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
		files = append(files, dat)
	}

	for _, file := range files {
		var data entities.Mappings
		err := json.Unmarshal(file, &data)
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

func GetMappingPath() ([]string, error) {
	var jsonPaths []string

	err := filepath.Walk("./mocks/mappings/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(info.Name()) == ".json" {
			jsonPaths = append(jsonPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return jsonPaths, nil
}
