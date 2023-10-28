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

type HintedHandoff struct {
	// TableName string                       `json:"TableName"`
	// Columns   []string                     `json:"Columns"`
	Rows map[NodeId][]AtomicDbMessage `json:"Row"`
}

type AtomicDbMessage struct {
	Data     Row    `json:"Data"`
	CourseId string `json:"CourseId"`
}

func (m AtomicDbMessage) String() string {
	return fmt.Sprintf("Data: %v, CourseId: %v", m.Data, m.CourseId)
}
func (h HintedHandoff) String() string {
	builder := &strings.Builder{}

	// Extract and sort the IDs
	ids := make([]string, 0, len(h.Rows))
	for id := range h.Rows {
		ids = append(ids, string(id))
	}
	sort.Strings(ids)

	// Print rows based on the sorted IDs
	fmt.Fprintln(builder, "Rows:")
	for _, id := range ids {
		messages := h.Rows[NodeId(id)]
		fmt.Fprintf(builder, "%s\n", id)
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

func (hhm *HintedHandoffManager) Append(nodeId NodeId, dbMsg AtomicDbMessage) error {
	hhm.mux.Lock()
	defer hhm.mux.Unlock()

	if hhm.Data.Rows == nil {
		hhm.Data.Rows = make(map[NodeId][]AtomicDbMessage)
	}

	hhm.Data.Rows[nodeId] = append(hhm.Data.Rows[nodeId], dbMsg)

	bytes, err := json.Marshal(hhm.Data)
	if err != nil {
		return err
	}

	return os.WriteFile(hhm.filepath, bytes, os.ModePerm)
}
