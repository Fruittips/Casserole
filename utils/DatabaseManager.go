package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Database, keyed by table name.
// Contains a set of columns defining the column key for each data item in each individual Row.
type Database struct {
	TableName    string           `json:"TableName"`
	PartitionKey string           `json:"PartitionKey"`
	Partitions   map[string][]Row `json:"Partitions"`
}

type Row struct {
	StudentId   string `json:"StudentId"`
	CreatedAt   int64  `json:"CreatedAt"`
	DeletedAt   int64  `json:"DeletedAt"`
	StudentName string `json:"StudentName"`
}

func (r Row) String() string {
	return fmt.Sprintf("StudentId: %s, CreatedAt: %d, DeletedAt: %d, StudentName: %s", r.StudentId, r.CreatedAt, r.DeletedAt, r.StudentName)
}

func (r Row) Equal(other Row) bool {
	return (r.StudentId == other.StudentId) &&
		(r.CreatedAt == other.CreatedAt) &&
		(r.DeletedAt == other.DeletedAt) &&
		(r.StudentName == other.StudentName)
}

func (d Database) String() string {
	builder := &strings.Builder{}

	// Print basic fields
	fmt.Fprintf(builder, "TableName: %s\n", d.TableName)

	// Extract and sort the partition keys
	partitionKeys := make([]string, 0, len(d.Partitions))
	for k := range d.Partitions {
		partitionKeys = append(partitionKeys, k)
	}
	sort.Strings(partitionKeys)
	fmt.Fprintf(builder, "=========\n")
	// Print rows based on the sorted partition keys
	fmt.Fprintln(builder, "Partitions:")
	fmt.Fprintf(builder, "=========\n")
	for _, partitionKey := range partitionKeys {
		rows := d.Partitions[partitionKey]
		fmt.Fprintf(builder, "%s: %s\n", d.PartitionKey, partitionKey)
		for _, row := range rows {
			fmt.Fprintf(builder, "\t\t%s\n", row)
		}
		fmt.Fprintf(builder, "------------\n")
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

func (db *DatabaseManager) AppendRow(partitionKey string, newData Row) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, exists := db.Data.Partitions[partitionKey]
	if !exists {
		db.Data.Partitions[partitionKey] = []Row{newData}
	} else {
		hasBeenUpdated := false
		for idx, pdata := range data { // handle write duplicates
			if pdata.StudentId == newData.StudentId {
				if newData.CreatedAt < pdata.CreatedAt { // assert that newData is later
					// return nil // possibly return as nil, data simply not added with no error handling
					return errors.New("write failed: earlier data found")
				} else {
					hasBeenUpdated = true
					data[idx] = newData
				}
			}
		}

		if !hasBeenUpdated {
			db.Data.Partitions[partitionKey] = append(data, newData)
		}
	}

	bytes, err := json.Marshal(db.Data)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.filepath, bytes, os.ModePerm)
	return err
}

func (db *DatabaseManager) GetRowByPartitionKey(courseId string, studentId string) (*Row, error) {
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
