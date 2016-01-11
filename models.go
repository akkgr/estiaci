package estia

import (
	"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Address struct {
	Area         string             `json:"area"`
	Street       string             `json:"street"`
	StreetNumber string             `json:"streetnumber"`
	PostalCode   string             `json:"postalcode"`
	Location     appengine.GeoPoint `json:"location"`
}

type Person struct {
	Lastname  string  `json:"lastname"`
	Firstname string  `json:"firstname"`
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
	Title    string `json:"title"`
	Position int32  `json:"position"`
	Resident Person `json:"resident"`
	Owner    Person `json:"owner"`
	Common   int64  `json:"common"`
	Elevetor int64  `json:"elevetor"`
	Heat     int64  `json:"heat"`
	Ei       int64  `json:"ei"`
	Fi       int64  `json:"fi"`
	Owners   int64  `json:"owners"`
	Other    int64  `json:"other"`
}

type Building struct {
	Id          int64        `json:"id" datastore:"-"`
	Address     Address      `json:"address"`
	Oil         int64        `json:"oil"`
	Fund        int64        `json:"fund"`
	Closed      int64        `json:"closed"`
	Active      bool         `json:"active"`
	Managment   bool         `json:"managment"`
	Appartments []Appartment `json:"appartments"`
}

type PublicBuild Building

func (b Building) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PublicBuild
		Title string `json:"title"`
	}{
		PublicBuild: PublicBuild(b),
		Title:       b.Address.Street + " " + b.Address.StreetNumber + ", " + b.Address.PostalCode + " " + b.Address.Area,
	})
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
