package estia

import (
	"net/http"
    
    "github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
    r.HandleFunc("/api/buildings", BuildAll).Methods("GET")
    r.HandleFunc("/api/buildings/{id}", BuildSingle).Methods("GET")
	r.HandleFunc("/api/buildings", BuildInsert).Methods("POST")
    r.HandleFunc("/api/buildings/{id}", BuildUpdate).Methods("PUT")
	r.HandleFunc("/api/buildings/{id}", BuildDelete).Methods("DELETE")

    handler := corsHandler(r)

    http.Handle("/", handler)
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
	} else {
			h.ServeHTTP(w, r)
		}
	}
}