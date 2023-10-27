// Commented so the system can run, this is still not converted to Ryan's version since he did the variant in func readRepair in read.go (line 94)

package utils

// import (
// 	"fmt"
// 	"sync"
// )

// type ReplicaData struct {
// 	filepath string
// 	data Database
// }

// type RowDiscrepancy struct {
// 	NodeId      NodeId
// 	CurrData    Row
// 	CorrectData Row
// }

// type ReadRepairsManager struct {
// 	mux           sync.Mutex
// 	Discrepancies []RowDiscrepancy
// 	Responses     []Response
// }

// //----------------------------------------
// // Print status??
// //----------------------------------------

// // func (rrm ReadRepairsManager) String() string {

// // }

// //----------------------------------------
// // Constructor
// //----------------------------------------

// func NewReadRepairsManager(responses []Response) *ReadRepairsManager {
// 	return &ReadRepairsManager{
// 		Discrepancies: make([]RowDiscrepancy, 0),
// 		Responses:     responses,
// 	}
// }

// //----------------------------------------
// // Methods
// //----------------------------------------
// func (rrm *ReadRepairsManager) PerformReadRepair(responses []Response) {
// 	// If there are no responses, there is nothing to repair
// 	if len(responses) == 0 {
// 		return
// 	}

// 	var latestData *Row
// 	var latestTimestamp int64
// 	var validResponses int

// 	// Identify the latest data
// 	for _, response := range responses {
// 		// If the response is not valid, skip it
// 		if response.Error != nil {
// 			// TODO : log error in response
// 			// TODO : sent to repair in for loop below, not sure if should
// 			continue
// 		}

// 		// handle or log nil data
// 		if response.Data == nil {
// 			// TODO : log error in response
// 			// TODO : sent to repair in for loop below, not sure if should
// 			continue
// 		}

// 		// If the response is valid, increment the number of valid responses
// 		validResponses++

// 		// If the response is the latest data, update the latest data
// 		if response.Data.CreatedAt > latestTimestamp {
// 			latestTimestamp = response.Data.CreatedAt
// 			latestData = response.Data
// 		}
// 	}

// 	// if there are no valid responses, there is nothing to repair
// 	if validResponses == 0 {
// 		return
// 	}

// 	// For each response, check if it matches the latest data
// 	for _, response := range responses {
// 		if response.Data.CreatedAt != latestTimestamp || response.Data == nil || response.Error != nil {
// 			discrepancy := RowDiscrepancy{
// 				NodeId:      response.NodeId,
// 				CurrData:    *response.Data,
// 				CorrectData: *latestData,
// 			}
// 			rrm.Discrepancies = append(rrm.Discrepancies, discrepancy)
// 		}
// 	}

// 	// Handle the discrepancies
// 	rrm.HandleDiscrepancies()
// }

// func (rrm *ReadRepairsManager) HandleDiscrepancies() {
// 	// If there are no discrepancies, there is nothing to repair
// 	if len(rrm.Discrepancies) == 0 {
// 		fmt.Println("[RRM] No discrepancies to repair")
// 		return
// 	}

// 	// Create a list to store write requests
// 	writeRequests := []Request{}

// 	// For each discrepancy, send a write request to the node with the latest data
// 	for _, discrepancy := range rrm.Discrepancies {
// 		writeURL := fmt.Sprintf(WRITE_ENDPOINT_FSTRING, discrepancy.NodeId, discrepancy.CorrectData.CourseId, discrepancy.CorrectData.StudentId)

// 		// Create a write request
// 		writeRequest := Request{
// 			NodeId:  discrepancy.NodeId,
// 			Url:     writeURL,
// 			Payload: &discrepancy.CorrectData,
// 		}

// 		// Add the write request to the list of write requests
// 		writeRequests = append(writeRequests, writeRequest)
// 	}

// 	// Send the write requests
// 	responses := IntraSystemRequests(writeRequests)

// 	// Handle the responses
// 	for _, response := range responses {
// 		if response.Error != nil {
// 			fmt.Println("[RRM] Error in repairing node ", response.NodeId, ": ", response.Error)
// 			// TODO : consider adding to hinted handoff or another mechanism to repair in the future
// 		} else {
// 			fmt.Println("[RRM] Successfully repaired node ", response.NodeId, " with data: ", response.Data)
// 		}
// 	}

// 	// Clear the discrepancies
// 	rrm.Discrepancies = []RowDiscrepancy{}
// }
