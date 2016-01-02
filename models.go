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

type Appartment struct {
	Title    string `json:"title" bson:"Title"`
	Position int    `json:"position" bson:"Position"`
}

type Building struct {
	Id          int64        `json:"id" datastore:"-"`
	Address     Address      `json:"address"`
	Oil         int64        `json:"oil"`
	Fund        int64        `json:"fund"`
	Active      bool         `json:"active"`
	Managment   bool         `json:"managment"`
	Appartments []Appartment `json:"appartments"`
}

type PublicBuild Building

func (b Building) MarshalJSON() ([]byte, error) {
	if b.Appartments == nil {
		b.Appartments = []Appartment{}
	}
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
