package utils

import (
	"sync"
)

// structs
type RepairData struct {
    PartitionKey string
    Entries      []StudentEntry
}

type StudentEntry struct {
    StudentId   string
    CreatedAt   int64
    DeletedAt   int64
    StudentName string
}

type ReadRepairsManager struct {
	filepath string
	mux      sync.Mutex
	Datas     []RepairData
}

func NewReadRepairsManager() *ReadRepairsManager {
	return &ReadRepairsManager{}
}


