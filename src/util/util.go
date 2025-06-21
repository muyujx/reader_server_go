package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

func RandomID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func MapToStruct(m map[string]any, out any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}
