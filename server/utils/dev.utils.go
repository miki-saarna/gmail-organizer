//go:build dev

package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error marshalling json: %v", err.Error())
	}
	return fmt.Sprintln(string(jsonData)), nil
}
