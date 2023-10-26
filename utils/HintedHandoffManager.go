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

type HintedHandoff struct {
	// TableName string                       `json:"TableName"`
	// Columns   []string                     `json:"Columns"`
	Rows map[string][]AtomicDbMessage `json:"Row"`
}

type AtomicDbMessage struct {
	Data      []Row `json:"Data"`
	Timestamp int64 `json:"Timestamp"`
}

func (m AtomicDbMessage) String() string {
	return fmt.Sprintf("Data: %v, Timestamp: %d", m.Data, m.Timestamp)
}
func (h HintedHandoff) String() string {
	builder := &strings.Builder{}

	// Print basic fields
	// fmt.Fprintf(builder, "TableName: %s\n", h.TableName)
	// fmt.Fprintf(builder, "Columns: %v\n", h.Columns)

	// Print rows with atomic messages
	fmt.Fprintln(builder, "Rows:")
	for id, messages := range h.Rows {
		fmt.Fprintf(builder, id, "\n")
		for _, msg := range messages {
			fmt.Fprintf(builder, "    - %s\n", msg)
		}
	}

	return builder.String()
}

type HintedHandoffManager struct {
	filepath string
	mux      sync.Mutex
	Data     HintedHandoff
}

func newHintedHandoffManager(path string) (*HintedHandoffManager, error) {
	if !filepath.IsAbs(path) {
		return nil, errors.New(fmt.Sprintf("Expected absolute path, was given %v", path))
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read file %v, error: %v", path, err))
	}

	var data HintedHandoff
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal JSON file %v, error: %v", path, err))
	}
	return &HintedHandoffManager{filepath: path, Data: data}, nil
}

func (hhm *HintedHandoffManager) Append(nodeId string, dbMsg AtomicDbMessage) error {
	hhm.mux.Lock()
	defer hhm.mux.Unlock()

	hhm.Data.Rows[nodeId] = append(hhm.Data.Rows[nodeId], dbMsg)

	bytes, err := json.Marshal(hhm.Data)
	if err != nil {
		return err
	}

	return os.WriteFile(hhm.filepath, bytes, os.ModePerm)
}
