package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/wesley601/mockid/db"
	"github.com/wesley601/mockid/handlers"
	"github.com/wesley601/mockid/services"
	"github.com/wesley601/mockid/utils"
)

func init() {
	if utils.Env("GO_ENV", "DEV") != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := db.EnsureDBFile(utils.Env("DB_PATH", "./data/local.db")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := app(utils.Env("DB_PATH", "./data/local.db")); err != nil {
		log.Fatalln(err)
	}
}

func app(dbPath string) error {
	conn, err := db.StartDB(dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()
	requestHandler := handlers.NewRequestHandler(conn)
	requestDAO := db.NewRequestDAO(conn)
	http.HandleFunc("GET /_/requests", requestHandler.Index)
	http.HandleFunc("GET /_/requests/{id}", requestHandler.Show)
	http.Handle("/", handlers.NewMapHandler(services.NewRequestMatcherLive(conn, requestDAO)))

	log.Println("Server starting at :3000")
	return http.ListenAndServe(":3000", nil)
}
