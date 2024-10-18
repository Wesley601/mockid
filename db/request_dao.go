package db

import (
	"database/sql"
)

type RequestDAO struct {
	db *sql.DB
}

func NewRequestDAO(db *sql.DB) *RequestDAO {
	return &RequestDAO{
		db: db,
	}
}

func (r *RequestDAO) Create(rr RequestSaved) error {
	_, err := r.db.Exec("INSERT INTO requests(requested_path, requested_method, matched_path, response_body, response_status) VALUES(?, ?, ?, ?, ?)",
		rr.RequestedPath, rr.RequestedMethod, rr.MatchedPath, rr.ResponseBody, rr.ResponseStatus)
	return err
}

func (r *RequestDAO) FindById(id string) (*RequestSaved, error) {
	var request RequestSaved
	row := r.db.QueryRow("SELECT id, requested_path, requested_method, matched_path, response_body, response_status, created_at FROM requests WHERE id=?;", id)
	if err := request.Scan(row); err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *RequestDAO) List() ([]RequestSaved, error) {
	var requests []RequestSaved
	rows, err := r.db.Query("SELECT id, requested_path, requested_method, matched_path, response_body, response_status, created_at FROM requests;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var request RequestSaved
		if err := request.Scan(rows); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *RequestDAO) Flush() error {
	_, err := r.db.Exec("DELETE FROM requests;")
	return err
}
