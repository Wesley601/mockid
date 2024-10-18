package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/wesley601/mockid/entities"
	"github.com/wesley601/mockid/services"
	"github.com/wesley601/mockid/views"
	"github.com/wesley601/mockid/views/mappings"
)

type MappingHandler struct{}

func NewMappingHandler() *MappingHandler {
	return &MappingHandler{}
}

func (re *MappingHandler) Show(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("filename")
	index, err := strconv.Atoi(r.PathValue("index"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data entities.Mappings
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, m := range data.Mappings {
		if i == index {
			m.FileName = path
			mappings.Show(m).Render(r.Context(), w)
			return
		}
	}
}

func (re *MappingHandler) Index(w http.ResponseWriter, r *http.Request) {
	var ms []entities.Mapping
	paths, err := services.GetMappingPath()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, m := range data.Mappings {
			m.FileName = path
			m.Index = i
			ms = append(ms, m)
		}
	}
	views.Base(mappings.Home(ms)).Render(r.Context(), w)
}
