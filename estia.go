package estia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var mySigningKey = []byte("TooSlowTooLate4u.")

func corsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.WriteHeader(http.StatusOK)
	} else {
		next(w, r)
	}
}

func authMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
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

	authBase := mux.NewRouter()
	router.PathPrefix("/auth").Handler(negroni.New(
		negroni.HandlerFunc(corsMiddleware),
		negroni.Wrap(authBase),
	))
	authRouter := authBase.PathPrefix("/auth").Subrouter()
	authRouter.Path("/login").Methods("POST").HandlerFunc(loginHandler)

	apiBase := mux.NewRouter()
	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(corsMiddleware),
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//validate user credentials
	if strings.ToLower(user.Username) != "admin" || user.Password != "123" {
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	// Create the Claims
	exp := time.Now().Add(time.Minute * 20).Unix()
	claims := &jwt.StandardClaims{
		ExpiresAt: exp,
		Issuer:    "test",
		Subject:   "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
	}

	//create a token instance using the token string
	response := Token{tokenString}
	jsonResponse(response, w)
}
