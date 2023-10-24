package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Database, keyed by table name.
// Contains a set of columns defining the column key for each data item in each individual Row.
type Database struct {
	TableName string      `json:"TableName"`
	Columns   []string    `json:"Columns"`
	Rows      map[int]Row `json:"Row"`
}

type Row struct {
	Data      []string `json:"Data"`
	Timestamp int64    `json:"Timestamp"`
}

func (r Row) String() string {
	return fmt.Sprintf("Data: %v, Timestamp: %d", r.Data, r.Timestamp)
}

func (d Database) String() string {
	builder := &strings.Builder{}

	// Print basic fields
	fmt.Fprintf(builder, "TableName: %s\n", d.TableName)
	fmt.Fprintf(builder, "Columns: %v\n", d.Columns)

	// Print rows
	fmt.Fprintln(builder, "Rows:")
	for id, row := range d.Rows {
		fmt.Fprintf(builder, "  ID %d -> %s\n", id, row)
	}

	return builder.String()
}

// Manages the database
type DatabaseManager struct {
	filepath string
	mux      sync.Mutex
	Data     Database
}

func newDatabaseManager(path string) *DatabaseManager {
	if !filepath.IsAbs(path) {
		panic(errors.New("the provided path is not an absolute path"))
	}
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var data Database
	json.Unmarshal(file, &data)

	return &DatabaseManager{filepath: path, Data: data}
}

func (db *DatabaseManager) AppendRow(newData Row) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	maxKey := 0
	for key := range db.Data.Rows {
		if key > maxKey {
			maxKey = key
		}
	}
	nextKey := maxKey + 1

	db.Data.Rows[nextKey] = newData

	bytes, err := json.Marshal(db.Data)
	if err != nil {
		return err
	}

	return os.WriteFile(db.filepath, bytes, os.ModePerm)
}

func (db *DatabaseManager) GetRowById(id int) (Row, error) {
	data, exists := db.Data.Rows[id]
	if !exists {
		return Row{}, errors.New("row does not exist")
	}
	return data, nil
}
