package services

import (
	"io/fs"
	"path/filepath"
)

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
