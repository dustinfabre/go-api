package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"net/http"
	"encoding/json"
)

type Config struct {
	Router 	*mux.Router
	DB		*sql.DB
}

func (c *Config) Conn(host string, port int, user string, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable",
	host, port, user, dbname)
	
	var err error
	c.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
	  panic(err)
	}

	err = c.DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Postgresql connected!")
}

func (c *Config) Start(port string) {
	// gorrilla router
	c.Routes()
    handler := cors.Default().Handler(c.Router)
	http.ListenAndServe(port, handler)
}

func (c *Config) Routes() {
	c.Router = mux.NewRouter()
	c.Router.HandleFunc("/v1/orders", c.getOrders).Methods("POST")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func (c *Config) getOrders(w http.ResponseWriter, r *http.Request) {
	var page Pagination
	json.NewDecoder(r.Body).Decode(&page)
	products, err := getOrders(c.DB, page.Start, page.End)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, products)
}