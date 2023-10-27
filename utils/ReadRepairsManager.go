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

//----------------------------------------
// Structs
//----------------------------------------

type Row struct {
	StudentId   string
	CreatedAt   int64
	DeletedAt   int64
	StudentName string
}

type Database struct {
	TableName    string           `json:"TableName"`
	PartitionKey int              `json:"PartitionKey"`
	Partitions   map[string][]Row `json:"Partitions"`
}

type ReplicaData struct {
	filepath string
	data     Database
}

type RowDiscrepancy struct {
	ReplicaData      []ReplicaData
	CurrentPartition int
	CurrData         []Row
	CorrectData      Row
}

type ReadRepairsManager struct {
	filepaths []string
	mux       sync.Mutex
	Datas     []ReplicaData
	NewDatas  []ReplicaData
}

//----------------------------------------
// Print status??
//----------------------------------------

func (rrm ReadRepairsManager) String() string {
	builder := &strings.Builder{}

	// Print basic fields (filepaths to compare)
	fmt.Fprintf(builder, "[RRM] Filepaths being compared: %v\n", rrm.filepaths)

	return builder.String()
}

//----------------------------------------
// Constructor
//----------------------------------------

func NewReadRepairsManager(filepaths []string) *ReadRepairsManager {
	rrm := ReadRepairsManager{
		filepaths: []string{},
		Datas:     []ReplicaData{},
	}

	// For each filepath given:
	for _, path := range filepaths {
		// validation : absolute path, readable file
		if !filepath.IsAbs(path) {
			panic(errors.New(fmt.Sprintf("Expected absolute path, was given %v", path)))
		}
		file, err := os.ReadFile(path)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Could not read file %v, error: %v", path, err)))
		}

		// unmarshal JSON file
		var replica ReplicaData
		err = json.Unmarshal(file, &replica.data)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Could not unmarshal JSON file %v, error: %v", path, err)))
		}

		// populate ReadRepairsManager's fields per replica
		rrm.filepaths = append(rrm.filepaths, path)
		rrm.Datas = append(rrm.Datas, replica)
	}

	fmt.Println("[RRM] Initialized ReadRepairsManager with filepaths: ", rrm.filepaths)

	return rrm, nil
}

//----------------------------------------
// Methods
//----------------------------------------
