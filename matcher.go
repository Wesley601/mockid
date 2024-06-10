package main

import (
	"database/sql"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Resp struct {
	Headers map[string]string
	Body    []byte
	Status  int
}

type RequestMatcherLive struct {
	db *sql.DB
}

func (m *RequestMatcherLive) Match(url, method string) (*Resp, error) {
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
		var data Mappings
		err := json.Unmarshal(file, &data)
		if err != nil {
			return nil, err
		}
		for _, mapping := range data.Mappings {
			if mapping.Request.Match(url, method) {
				response, err := mapping.Response.GetBody()
				if err != nil {
					return nil, err
				}
				_, err = m.db.Exec("INSERT INTO requests(requested_path, requested_method, matched_path, response_body, response_status) VALUES(?, ?, ?, ?, ?)",
					url, method, mapping.Request.URLPattern, string(response), mapping.Response.Status)
				if err != nil {
					log.Printf("Unable to insert on requests table, %s\n", err.Error())
				}
				return &Resp{
					Status:  mapping.Response.Status,
					Headers: mapping.Response.Headers,
					Body:    response,
				}, nil
			}
		}
	}

	return nil, errRespNotFound
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
