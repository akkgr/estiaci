package estia

import (
	"net/http"
)

func init() {
	router := NewRouter()
	http.Handle("/", corsHandler(router))
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
