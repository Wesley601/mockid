package db

import (
	"fmt"
)

type RequestSaved struct {
	ID              int
	RequestedPath   string
	RequestedMethod string
	MatchedPath     string
	ResponseBody    string
	ResponseStatus  int
	CreatedAt       string
}

type Scanner interface {
	Scan(dest ...any) error
}

func (r *RequestSaved) Scan(row Scanner) error {
	if err := row.Scan(&r.ID, &r.RequestedPath, &r.RequestedMethod, &r.MatchedPath, &r.ResponseBody, &r.ResponseStatus, &r.CreatedAt); err != nil {
		return fmt.Errorf("failed to scan a request row: %w", err)
	}
	return nil
}
