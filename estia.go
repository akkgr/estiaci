package estia

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func corsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Offset,Limit")
		w.WriteHeader(http.StatusOK)
	} else {
		next(w, r)
	}
}

func authMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r,  request.OAuth2Extractor ,func(token *jwt.Token) (interface{}, error) {
		return "secret", nil
	})

	if err == nil {

		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}
}

func init() {

	r := mux.NewRouter()
	n := negroni.New(negroni.HandlerFunc(corsMiddleware), negroni.HandlerFunc(authMiddleware))

	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/buildings", BuildAll).Methods("GET")
	r.HandleFunc("/api/buildings/{id}", BuildSingle).Methods("GET")
	r.HandleFunc("/api/buildings", BuildInsert).Methods("POST")
	r.HandleFunc("/api/buildings/{id}", BuildUpdate).Methods("PUT")
	r.HandleFunc("/api/buildings/{id}", BuildDelete).Methods("DELETE")

	n.UseHandler(r)

	http.Handle("/", n)
}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err :=  json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
