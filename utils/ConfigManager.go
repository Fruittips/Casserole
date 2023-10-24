package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manages the configuration, based on the given config filepath.
type ConfigManager struct {
	filepath         string
	defaultNodes     map[int]Node
	ConsistencyLevel string // e.g. "QUORUM"
	GracePeriod      time.Duration
	Timeout          time.Duration
	RF               int // replication factor
}

// Load configuration from a given path.
func newConfigManager(path string) *ConfigManager {
	if !filepath.IsAbs(path) {
		panic(errors.New("the provided path is not an absolute path"))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading config JSON: %v", err))
	}

	var data config
	err = json.Unmarshal(file, &data)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling config JSON: %v", err))
	}

	// Convert hardcoded node data into Nodes
	nodeMap := make(map[int]Node, 0)
	for id, configNode := range data.Ring {
		nodeMap[id] = Node{
			Id:     id,
			Port:   configNode.Port,
			IsDead: configNode.IsDead,
		}
	}

	return &ConfigManager{
		path,
		nodeMap,
		data.ConsistencyLevel,
		data.GracePeriod,
		data.Timeout,
		data.RF,
	}
}

func (cMgr *ConfigManager) String() string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "ConsistencyLevel: %s\n", cMgr.ConsistencyLevel)
	fmt.Fprintf(builder, "GracePeriod: %s\n", cMgr.GracePeriod)
	fmt.Fprintf(builder, "Timeout: %s\n", cMgr.Timeout)
	fmt.Fprintf(builder, "RF: %d\n", cMgr.RF)

	fmt.Fprintln(builder, "Ring:")
	for id, node := range cMgr.defaultNodes {
		fmt.Fprintf(builder, "  ID %d -> %s\n", id, node)
	}

	return builder.String()
}

type config struct {
	ConsistencyLevel string             `json:"consistencyLevel"` //e.g. "QUORUM"
	GracePeriod      time.Duration      `json:"gracePeriod" `     //duration in seconds
	Timeout          time.Duration      `json:"timeout"`          //duration in seconds
	RF               int                `json:"rf"`               //replication factor
	Ring             map[int]configNode `json:"ring"`             //all nodes in ring
}

type configNode struct {
	Port   int  `json:"port"`
	IsDead bool `json:"isDead"`
}

func (n configNode) String() string {
	deadStatus := "Alive"
	if n.IsDead {
		deadStatus = "Dead"
	}
	return fmt.Sprintf("Port: %d, Status: %s", n.Port, deadStatus)
}
