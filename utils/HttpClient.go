/**
Defines HTTP interactions between nodes.
*/

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Sends an internal read to the given node. Returns a data Row or an error.
func (nm *NodeManager) SendInternalRead(dstNode Node, courseId string, studentId string) (*Row, error) {
	timeout := time.Duration(nm.GetConfig().Timeout)
	client := &http.Client{Timeout: timeout * time.Second} //TODO: should second be already in there?

	// Generate the URL of the node: Target port, Read-from CourseID, Read-from StudentID
	url := fmt.Sprintf(BASE_URL+INTERNAL_READ_ENDPOINT_FSTRING, dstNode.Port, courseId, studentId)

	// Send the GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("HTTP Error: %v", resp.Status))
	}

	// Parse the response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert response body into HTTPResponse, then into Row
	var row Row
	err = json.Unmarshal(body, &row)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

// Sends an internal write to the given node.
func (nm *NodeManager) SendInternalWrite(dstNode Node, courseId string, data Row) error {
	timeout := time.Duration(nm.GetConfig().Timeout)
	client := &http.Client{Timeout: timeout * time.Second} //TODO: should second be already in there?

	// Generate the URL of the node: Target port, Write-to CourseID, Write-to StudentID
	url := fmt.Sprintf(BASE_URL+INTERNAL_WRITE_ENDPOINT_FSTRING, dstNode.Port, courseId)

	// Send the POST request
	newRowJSON, err := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(newRowJSON))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("HTTP Error: %v", resp.Status))
	}

	return nil
}

// Sends an internal hinted handoffs request to the given node. This triggers the given node to send any necessary HHs with internal write requests.
func (nm *NodeManager) RequestForHHs(dstNode Node) error {
	timeout := time.Duration(nm.GetConfig().Timeout)
	client := &http.Client{Timeout: timeout * time.Second} //TODO: should second be already in there?

	// Generate the URL of the node: Target port, my own ID (so target knows who wants the hhs)
	url := fmt.Sprintf(BASE_URL+INTERNAL_CHECKHH_ENDPOINT_FSTRING, dstNode.Port, nm.LocalId)

	// Send the GET request (TODO: if handler is modified, this may need to change too)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("HTTP Error: %v", resp.Status))
	}

	return nil
}
