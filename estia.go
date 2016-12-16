package estia

import (
	"net/http"

	"github.com/rs/cors"
    "github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
    r.HandleFunc("/api/buildings", BuildAll).Methods("GET")
    r.HandleFunc("/api/buildings/{id}", BuildSingle).Methods("GET")
	r.HandleFunc("/api/buildings", BuildInsert).Methods("POST")
    r.HandleFunc("/api/buildings/{id}", BuildUpdate).Methods("PUT")
	r.HandleFunc("/api/buildings/{id}", BuildDelete).Methods("DELETE")

    handler := cors.Default().Handler(r)

    http.Handle("/", handler)
}
