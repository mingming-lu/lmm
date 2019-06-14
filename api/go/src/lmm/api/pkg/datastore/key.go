package datastore

import "cloud.google.com/go/datastore"

func MustKey(encoded string) *datastore.Key {
	key, err := datastore.DecodeKey(encoded)
	if err != nil {
		panic(err)
	}
	return key
}
