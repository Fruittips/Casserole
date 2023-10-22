package handlers

import (
	"encoding/json"
	"fmt"
	"os"
)

type AtomicDbMessage struct {
	Data      []string `json: "Data"`
	Timestamp int64    `json: "Timestamp"`
}

type Data struct {
	TableName string                     `json: "TableName"`
	Columns   []string                   `json: "Columns"`
	Row       map[string]AtomicDbMessage `json: "Row"`
}

func internalRead(toNode int) (Data, error) {
	var nodeData Data
	// use err to determine if there is an existing json file
	filename := fmt.Sprintf("dbFiles/node-%d.json", toNode)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		emptyData := Data{TableName: "", Columns: make([]string, 0), Row: make(map[string]AtomicDbMessage)}
		return emptyData, err
	}

	// parse ogData
	byteValue, err := os.ReadFile(filename)
	check(err)
	_ = json.Unmarshal([]byte(byteValue), &nodeData)

	return nodeData, nil
}
