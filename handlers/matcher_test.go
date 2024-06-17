package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wesley601/mockid/db"
	"github.com/wesley601/mockid/entities"
	"github.com/wesley601/mockid/handlers"
	"github.com/wesley601/mockid/services"
	"github.com/wesley601/mockid/utils"
)

func init() {
	utils.SetToRoot("..")
}

func Test_Matcher(t *testing.T) {
	conn, err := db.StartDB(":memory:")
	if err != nil {
		panic(err)
	}
	requestDAO := db.NewRequestDAO(conn)
	handler := handlers.NewMapHandler(services.NewRequestMatcherLive(conn, requestDAO))

	testCases := []struct {
		request  *http.Request
		response struct {
			filePath string
			code     int
		}
		desc string
	}{
		{
			request: utils.Must(http.NewRequest(http.MethodGet, "/accounts/1", nil)),
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
			request: utils.Must(http.NewRequest(http.MethodGet, "/v1/card?externalId=A10000301&product=stub", nil)),
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
			want, err := entities.GetBody(tC.response.filePath)
			if err != nil {
				t.Errorf("unable to get the file response from %q\n err: %q", tC.response, err.Error())
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
	conn, err := db.StartDB(":memory:")
	if err != nil {
		panic(err)
	}
	requestDAO := db.NewRequestDAO(conn)

	handler := handlers.NewMapHandler(services.NewRequestMatcherLive(conn, requestDAO))

	t.Run("should save the request in the database", func(t *testing.T) {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, utils.Must(http.NewRequest(http.MethodGet, "/accounts/1", nil)))

		requests, err := requestDAO.List()
		if err != nil {
			t.Errorf("something went wrong %q", err.Error())
			return
		}
		if len(requests) == 0 {
			t.Errorf("the request was not saved on the database")
		}
	})
}
