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

    // methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
    // headers := handlers.AllowedHeaders([]string{"*"})
    // origins := handlers.AllowedOrigins([]string{"http://localhost/"})
    // options := handlers.IgnoreOptions()
    handler := cors.Default().Handler(r)

    http.Handle("/", handler)
}
