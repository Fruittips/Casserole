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
	TableName string                    `json:"TableName"`
	Columns   []string                  `json:"Columns"`
	Rows      map[int][]AtomicDbMessage `json:"Row"`
}

type AtomicDbMessage struct {
	Data      []string `json:"Data"`
	Timestamp int64    `json:"Timestamp"`
}

func (m AtomicDbMessage) String() string {
	return fmt.Sprintf("Data: %v, Timestamp: %d", m.Data, m.Timestamp)
}
func (h HintedHandoff) String() string {
	builder := &strings.Builder{}

	// Print basic fields
	fmt.Fprintf(builder, "TableName: %s\n", h.TableName)
	fmt.Fprintf(builder, "Columns: %v\n", h.Columns)

	// Print rows with atomic messages
	fmt.Fprintln(builder, "Rows:")
	for id, messages := range h.Rows {
		fmt.Fprintf(builder, "  ID %d:\n", id)
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

func newHintedHandoffManager(path string) *HintedHandoffManager {
	if !filepath.IsAbs(path) {
		panic(errors.New("the provided path is not an absolute path"))
	}
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var data HintedHandoff
	json.Unmarshal(file, &data)

	return &HintedHandoffManager{filepath: path, Data: data}
}

func (hhm *HintedHandoffManager) Append(nodeId int, dbMsg AtomicDbMessage) error {
	hhm.mux.Lock()
	defer hhm.mux.Unlock()

	hhm.Data.Rows[nodeId] = append(hhm.Data.Rows[nodeId], dbMsg)

	bytes, err := json.Marshal(hhm.Data)
	if err != nil {
		return err
	}

	return os.WriteFile(hhm.filepath, bytes, os.ModePerm)
}
