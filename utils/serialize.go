package utils

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
)

func Serialize(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func Deserialize(data []byte, out interface{}) error {
	if len(data) == 0 {
		return errors.New("no data to deserialize")
	}

	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}

func DecodeMessage(obj any, value []byte) error {
	if err := json.Unmarshal(value, &obj); err != nil {
		log.Printf("Error: Failed to decode message: %s", err.Error())
		return errors.New("error: failed to decode message")
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		log.Printf("Error: Failed to validate message: %s", err.Error())
		return errors.New("error: failed to validate message")
	}

	return nil
}
