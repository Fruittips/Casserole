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

type Database struct {
	TableName    string        `json:"TableName"`
	PartitionKey int           `json:"PartitionKey"`
	Parititons   map[int][]Row `json:"Parititons"`
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
	for partitionKey, rows := range d.Parititons {
		fmt.Fprintf(builder, "\tPartitionKey: %d\n", partitionKey)
		for _, row := range rows {
			fmt.Fprintf(builder, "\t\t%s\n", row)
		}
	}

	return builder.String()
}

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

func (db *DatabaseManager) AppendRow(partitionKey int, newData Row) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, exists := db.Data.Parititons[partitionKey]
	if !exists {
		db.Data.Parititons[partitionKey] = []Row{newData}
	} else {
		db.Data.Parititons[partitionKey] = append(data, newData)
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
	data, exists := db.Data.Parititons[courseId]
	if !exists {
		return nil, errors.New("Parititon not found")
	}

	for _, row := range data {
		if row.StudentId == studentId {
			return &row, nil
		}
	}

	return nil, errors.New("Row not found")

}
