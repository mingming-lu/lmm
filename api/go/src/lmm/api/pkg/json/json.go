package json

import "encoding/json"

// JSON is a shortcut for map[string]interface{}
type JSON map[string]interface{}

// JSONify encode obj into string
func JSONify(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// MustJSONify panics if obj cannot be jsonified
func MustJSONify(obj interface{}) string {
	s, err := JSONify(obj)
	if err != nil {
		panic(err)
	}
	return s
}
