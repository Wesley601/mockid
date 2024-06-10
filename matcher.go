package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type RequestMatcherLive struct {
	Paths []string
}

func NewResquetMatcherLive() (*RequestMatcherLive, error) {
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

	return &RequestMatcherLive{
		Paths: jsonPaths,
	}, nil
}

func (m *RequestMatcherLive) Match(url, method string) (*Resp, error) {
	var files [][]byte
	for _, v := range m.Paths {
		dat, err := os.ReadFile(v)
		if err != nil {
			log.Printf("Unable to read the file %s: %s", v, err.Error())
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
