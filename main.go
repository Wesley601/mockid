package main

import (
	"errors"
	"log"
	"net/http"
)

func init() {
	err := EnsureDBFile("./data/local.db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := app("./data/local.db"); err != nil {
		log.Fatalln(err)
	}
}

func app(dbPath string) error {
	db, err := StartDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	http.HandleFunc("GET /_/requests", func(w http.ResponseWriter, r *http.Request) {
		var requests []RequestSaved
		rows, err := db.Query("SELECT id, requested_path, requested_method, matched_path, response_body, response_status FROM requests;")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var request RequestSaved
			if err := request.Scan(rows); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			requests = append(requests, request)
		}
		Base(Home(requests)).Render(r.Context(), w)
	})

	http.HandleFunc("GET /_/requests/{id}", func(w http.ResponseWriter, r *http.Request) {
		var request RequestSaved
		row := db.QueryRow("SELECT id, requested_path, requested_method, matched_path, response_body, response_status FROM requests WHERE id=?", r.PathValue("id"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := request.Scan(row); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Show(request).Render(r.Context(), w)
	})

	http.Handle("/", &MapHandler{matcher: &RequestMatcherLive{db: db}})

	log.Println("Server starting at :3000")
	return http.ListenAndServe(":3000", nil)
}

var errRespNotFound = errors.New("response not found")

type RequestMatcher interface {
	Match(r *http.Request) (*Resp, error)
}

type MapHandler struct {
	matcher RequestMatcher
}

func (m *MapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := m.matcher.Match(r)
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
