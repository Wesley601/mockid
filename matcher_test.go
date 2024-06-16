package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Matcher(t *testing.T) {
	db, err := StartDB(":memory:")
	if err != nil {
		panic(err)
	}
	handler := &MapHandler{matcher: &RequestMatcherLive{db: db}}

	testCases := []struct {
		request  *http.Request
		response struct {
			filePath string
			code     int
		}
		desc string
	}{
		{
			request: Must(http.NewRequest(http.MethodGet, "/accounts/1", nil)),
			response: struct {
				filePath string
				code     int
			}{
				filePath: "account-service/getAccount.json",
				code:     200,
			},
			desc: "should match the url",
		},
		{
			request: Must(http.NewRequest(http.MethodGet, "/v1/card?externalId=A10000301&product=stub", nil)),
			response: struct {
				filePath string
				code     int
			}{
				filePath: "pismo-service/getCardProducts.json",
				code:     200,
			},
			desc: "should match the url with query params",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, tC.request)
			got := response.Body.String()
			want, err := GetBody(tC.response.filePath)
			if err != nil {
				t.Errorf("unable to get the file response from %q", tC.response)
				return
			}

			if got != string(want) {
				t.Errorf("invalid body got %q, want %q", got, want)
			}

			if response.Result().StatusCode != response.Code {
				t.Errorf("invalid http code got %q, want %q", response.Result().StatusCode, response.Code)
			}
		})
	}
}

func Test_SaveMatcher(t *testing.T) {
	db, err := StartDB(":memory:")
	if err != nil {
		panic(err)
	}
	handler := &MapHandler{matcher: &RequestMatcherLive{db: db}}

	t.Run("should save the request in the database", func(t *testing.T) {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, Must(http.NewRequest(http.MethodGet, "/accounts/1", nil)))

		var request RequestSaved
		row := db.QueryRow("SELECT id, requested_path, requested_method, matched_path, response_body, response_status FROM requests;")

		if request.Scan(row) != nil {
			if errors.Is(err, sql.ErrNoRows) {
				t.Errorf("the request was not saved on the database")
			} else {
				t.Errorf("something went wrong %q", err.Error())
			}
		}
	})
}
