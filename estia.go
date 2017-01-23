package estia

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var mySigningKey = []byte("secret")

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
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
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
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Path("/login").Methods("POST").HandlerFunc(LoginHandler)

	apiBase := mux.NewRouter()
	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(authMiddleware),
		negroni.Wrap(apiBase),
	))
	apiRouter := apiBase.PathPrefix("/api").Subrouter()
	apiRouter.Path("/buildings").Methods("GET").HandlerFunc(BuildAll)
	apiRouter.Path("/buildings/{id}").Methods("GET").HandlerFunc(BuildSingle)
	apiRouter.Path("/buildings").Methods("POST").HandlerFunc(BuildInsert)
	apiRouter.Path("/buildings/{id}").Methods("PUT").HandlerFunc(BuildUpdate)
	apiRouter.Path("/buildings/{id}").Methods("DELETE").HandlerFunc(BuildDelete)

	http.Handle("/", router)
}

func jsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
