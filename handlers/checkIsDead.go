package handlers

import (
	"encoding/json"
	"os"
)

func checkIsDead() (bool, error) {
	// TODO: Remove, shift to NodeManager -- also, how does this know the ID of the node to check?
	// check if config.json is_dead is false
	filename := "config.json"

	var filestruc map[string]interface{}

	// parse ogData
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return true, err
	}
	_ = json.Unmarshal([]byte(byteValue), &filestruc)

	if filestruc["is_dead"] == true {
		return true, nil
	}
	return false, nil
}
