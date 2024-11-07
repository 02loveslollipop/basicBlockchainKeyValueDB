package blockchain
import (
	"encoding/json"
)

// key-value abstraction in the blockchain
type KeyValue struct {
	Key   string `json:"key"` //Key of the key-value pair
	Value string `json:"value"` //Value of the key-value pair
}

// Key-Value functions
func ToJSON(keyValue KeyValue) ([]byte, error) {
	return json.Marshal(keyValue)
}

func FromJSON(data []byte) (KeyValue, error) {
	var keyValue KeyValue
	err := json.Unmarshal(data, &keyValue)
	if err != nil {
		return KeyValue{}, err
	}
	return keyValue, nil
}

func FromString(data string) (KeyValue, error) {
	return FromJSON([]byte(data))
}

func ToString(keyValue KeyValue) (string, bool) {
	data, err := ToJSON(keyValue)
	if err != nil {
		return "", false
	}
	return string(data), true
}