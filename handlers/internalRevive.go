package handlers

import (
	"encoding/json"
	"os"
)

func internalRevive() error {
	// check if config.json is_dead is false
	filename := "config.json"

	var filestruc map[string]interface{}

	// parse ogData
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	_ = json.Unmarshal([]byte(byteValue), &filestruc)

	filestruc["is_dead"] = false

	byteFile, err3 := json.MarshalIndent(filestruc, "", "\t")
	if err3 != nil {
		return err3
	}

	// write byte into file
	err4 := os.WriteFile(filename, byteFile, 0644)
	if err3 != nil {
		return err4
	}

	// let all files know that it exists by talking to all ports
	// TODO: requires ports to finaliez first

	return nil

}
