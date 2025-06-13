package common

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
)

// SelectDataType is a simple key/value pair that can be used to populate
// <select> dropdowns or any data-driven lists in templates.
type SelectDataType struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

// ToPointer returns a pointer to the provided value.
func ToPointer[T any](data T) *T { return &data }

const (
	LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberBytes = "0123456789"
)

// GenerateToken produces a random alpha-numeric string of the given length.
func GenerateToken(n int) string {
	allBytes := LetterBytes + NumberBytes
	b := make([]byte, n)
	for i := range b {
		b[i] = allBytes[rand.Intn(len(allBytes))]
	}
	return string(b)
}

// Print pretty-prints any value as indented JSON to stdout.
func Print(data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}

// ToSelect converts a slice of structs into a SelectDataType slice based on field names.
func ToSelect[T any](data []T, keyField, valField string) ([]SelectDataType, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("data should be a slice")
	}

	var result []SelectDataType
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		k := reflect.Indirect(reflect.ValueOf(item)).FieldByName(keyField)
		val := reflect.Indirect(reflect.ValueOf(item)).FieldByName(valField)
		if !k.IsValid() || !val.IsValid() {
			return nil, fmt.Errorf("invalid field names")
		}
		result = append(result, SelectDataType{
			Key: fmt.Sprintf("%v", k.Interface()),
			Val: fmt.Sprintf("%v", val.Interface()),
		})
	}
	return result, nil
}
