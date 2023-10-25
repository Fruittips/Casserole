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
	TableName    string        `json:"TableName"`
	PartitionKey int           `json:"PartitionKey"`
	Partitions   map[int][]Row `json:"Partitions"`
}

type Row struct {
	StudentId   int    `json:"StudentId"`
	CreatedAt   int64  `json:"CreatedAt"`
	DeletedAt   int64  `json:"DeletedAt"`
	StudentName string `json:"StudentName"`
}

func (r Row) String() string {
	return fmt.Sprintf("StudentId: %d, CreatedAt: %d, DeletedAt: %d, StudentName: %s", r.StudentId, r.CreatedAt, r.DeletedAt, r.StudentName)
}

func (d Database) String() string {
	builder := &strings.Builder{}

	// Print basic fields
	fmt.Fprintf(builder, "TableName: %s\n", d.TableName)

	// Print rows
	fmt.Fprintln(builder, "Partitions:")
	for partitionKey, rows := range d.Partitions {
		fmt.Fprintf(builder, "\tPartitionKey: %d\n", partitionKey)
		for _, row := range rows {
			fmt.Fprintf(builder, "\t\t%s\n", row)
		}
	}

	return builder.String()
}

// Manages the database
type DatabaseManager struct {
	filepath string
	mux      sync.Mutex
	Data     Database
}

func newDatabaseManager(path string) (*DatabaseManager, error) {
	if !filepath.IsAbs(path) {
		return nil, errors.New(fmt.Sprintf("Expected absolute path, was given %v", path))
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read file %v, error: %v", path, err))
	}

	var data Database
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal JSON file %v, error: %v", path, err))
	}

	return &DatabaseManager{filepath: path, Data: data}, nil
}

func (db *DatabaseManager) AppendRow(partitionKey int, newData Row) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, exists := db.Data.Partitions[partitionKey]
	if !exists {
		db.Data.Partitions[partitionKey] = []Row{newData}
	} else {
		db.Data.Partitions[partitionKey] = append(data, newData)
	}

	bytes, err := json.Marshal(db.Data)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.filepath, bytes, os.ModePerm)
	return err
}

func (db *DatabaseManager) GetRowByPartitionKey(courseId int, studentId int) (*Row, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, exists := db.Data.Partitions[courseId]
	if !exists {
		return nil, errors.New("Partition not found")
	}

	for _, row := range data {
		if row.StudentId == studentId {
			return &row, nil
		}
	}

	return nil, errors.New("Row not found")

}
