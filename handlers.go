package estia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	mySigningKey := []byte("secret")

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
	JsonResponse(response, w)
}

//BuildAll List of all Buildings
func BuildAll(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if u := user.Current(c); u != nil {
	}

	offset := r.Header.Get("Offset")
	limit := r.Header.Get("Limit")

	q := datastore.NewQuery("Building")
	q = q.Order("Address.Street")
	q = q.Order("Address.StreetNumber")
	q = q.Order("Address.Area")

	i, err := strconv.ParseInt(offset, 10, 32)
	if err == nil {
		q = q.Offset(int(i))
	}
	i, err = strconv.ParseInt(limit, 10, 32)
	if err == nil {
		q = q.Limit(int(i))
	}

	result := []Building{}

	keys, err := q.GetAll(c, &result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(result); i++ {
		result[i].Id = keys[i].IntID()
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

//BuildSingle Building by id
func BuildSingle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result Building
	result.Id = i

	if i != 0 {
		c := appengine.NewContext(r)
		k := result.key(c)
		err = datastore.Get(c, k, &result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result.Id = k.IntID()
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//BuildInsert insert a Building
func BuildInsert(w http.ResponseWriter, r *http.Request) {
	build := Building{}
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if build.Appartments == nil {
		build.Appartments = []Appartment{{Title: "A", Position: 1}}
	}

	c := appengine.NewContext(r)

	err := build.save(c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//BuildUpdate Update a Building
func BuildUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	build := Building{}
	build.Id = i
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)

	if build.Appartments == nil {
		build.Appartments = []Appartment{{Title: "A", Position: 1}}
	}

	err = build.save(c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(build); err != nil {
		panic(err)
	}
}

//BuildDelete Delete a Building
func BuildDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)

	build := Building{}
	build.Id = i
	err = datastore.Delete(c, build.key(c))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
