package estia

import (
	"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type GeoLocation appengine.GeoPoint

type Address struct {
	Area         string      `json:"area"`
	Street       string      `json:"street"`
	StreetNumber string      `json:"streetNumber"`
	PostalCode   string      `json:"postalCode"`
	Country      string      `json:"country"`
	Location     GeoLocation `json:"location"`
}

type Person struct {
	Id        int64   `json:"id" datastore:"-"`
	Lastname  string  `json:"lastname"`
	Firstname string  `json:"firstname"`
	Fullname  string  `json:"fullname" datastore:"-"`
	Address   Address `json:"address"`
	Home      string  `json:"home"`
	Work      string  `json:"work"`
	Mobile    string  `json:"mobile"`
	Fax       string  `json:"fax"`
	Other     string  `json:"other"`
	Email     string  `json:"email"`
	Ibank     string  `json:"ibank"`
}

type Appartment struct {
	Id       int64    `json:"id" datastore:"-"`
	Title    string   `json:"title"`
	Position int32    `json:"position"`
	Resident Person   `json:"resident"`
	Owner    Person   `json:"owner"`
	Common   int64    `json:"common"`
	Elevetor int64    `json:"elevetor"`
	Heat     int64    `json:"heat"`
	Ei       int64    `json:"ei"`
	Fi       int64    `json:"fi"`
	Owners   int64    `json:"owners"`
	Other    int64    `json:"other"`
	Counters []string `json:"counters"`
}

type Building struct {
	Id          int64        `json:"id" datastore:"-"`
	Address     Address      `json:"address"`
	Oil         int64        `json:"oil"`
	Fund        int64        `json:"fund"`
	Closed      int64        `json:"closed"`
	Active      bool         `json:"active"`
	Managment   bool         `json:"managment"`
	Appartments []Appartment `json:"appartments" datastore:"-"`
}

func (b *Building) key(c context.Context) *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if b.Id == 0 {
		return datastore.NewIncompleteKey(c, "Building", nil)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	return datastore.NewKey(c, "Building", "", b.Id, nil)
}

func (b *Building) save(c context.Context) error {
	// reference the key function and generate it
	// accordingly basically its isNew true/false
	k, err := datastore.Put(c, b.key(c), b)
	if err != nil {
		return err
	}

	// The Id on the model is not prepopulated so we'll have
	// to append manually
	b.Id = k.IntID()
	return nil
}

func (g *GeoLocation) UnmarshalJSON(b []byte) (err error) {
	var jm map[string]float64
	if err = json.Unmarshal(b, &jm); err == nil {
		g.Lat = jm["lat"]
		g.Lng = jm["lng"]
	}
	return
}

func (g GeoLocation) MarshalJSON() ([]byte, error) {
	jm := make(map[string]float64)
	jm["lat"] = g.Lat
	jm["lng"] = g.Lng
	return json.Marshal(jm)
}
