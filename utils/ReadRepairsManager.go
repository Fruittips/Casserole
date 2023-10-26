package utils

import {
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
}

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


