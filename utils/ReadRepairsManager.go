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
	// If there are no responses, there is nothing to repair
	if len(responses) == 0 {
		return
	}

	var latestData *Row
	var latestTimestamp int64
	var validResponses int

	// Identify the latest data
	for _, response := range responses {
		// If the response is not valid, skip it
		if response.Error != nil {
			// TODO : log error in response
			// TODO : sent to repair in for loop below, not sure if should
			continue
		}

		// handle or log nil data
		if response.Data == nil {
			// TODO : log error in response
			// TODO : sent to repair in for loop below, not sure if should
			continue
		}

		// If the response is valid, increment the number of valid responses
		validResponses++

		// If the response is the latest data, update the latest data
		if response.Data.CreatedAt > latestTimestamp {
			latestTimestamp = response.Data.CreatedAt
			latestData = response.Data
		}
	}

	// if there are no valid responses, there is nothing to repair
	if validResponses == 0 {
		return
	}

	// For each response, check if it matches the latest data
	for _, response := range responses {
		if response.Data.CreatedAt != latestTimestamp || response.Data == nil || response.Error != nil {
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

func (rrm *ReadRepairsManager) HandleDiscrepancies() {
	// For each discrepancy, send a write request to the node with the latest data
	for _, discrepancy := range rrm.Discrepancies {
		// TODO : send write request to node with latest data

	}

	// Clear the discrepancies
	rrm.Discrepancies = []RowDiscrepancy{}
}
