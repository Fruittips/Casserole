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

type RowDiscrepancy struct {
	NodeId      NodeId
	CurrData    Row
	CorrectData Row
}

type ReadRepairsManager struct {
	mux           sync.Mutex
	Discrepancies []RowDiscrepancy
}

//----------------------------------------
// Print status??
//----------------------------------------

// func (rrm ReadRepairsManager) String() string {

// }

//----------------------------------------
// Constructor
//----------------------------------------

func NewReadRepairsManager(filepaths []string, responses []Response) *ReadRepairsManager {
	return &ReadRepairsManager{
		Discrepancies: make([]RowDiscrepancy, 0),
	}
}

//----------------------------------------
// Methods
//----------------------------------------
func (rrm *ReadRepairsManager) PerformReadRepair(responses []Response) {
	var latestData *Row
	var latestTimestamp int64

	// Identify the latest data
	for _, response := range responses {
		if response.Data.CreatedAt > latestTimestamp {
			latestTimestamp = response.Data.CreatedAt
			latestData = response.Data
		}
	}

	// For each response, check if it matches the latest data
	for _, response := range responses {
		if response.Data.CreatedAt != latestTimestamp {
			discrepancy := RowDiscrepancy{
				NodeId:      response.NodeId,
				CurrData:    *response.Data,
				CorrectData: *latestData,
			}
			rrm.Discrepancies = append(rrm.Discrepancies, discrepancy)
		}
	}

	// Handle the discrepancies
	rrm.HandleDiscrepancies()
}
