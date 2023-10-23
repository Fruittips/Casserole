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

type Config struct {
	ConsistencyLevel string        `json:"consistencyLevel"` //e.g. "QUORUM"
	GracePeriod      time.Duration `json:"gracePeriod" `     //duration in seconds
	Timeout          time.Duration `json:"timeout"`          //duration in seconds
	RF               int           `json:"rf"`               //replication factor
	Ring             map[int]Node  `json:"ring"`             //all nodes in ring
}

func (c Config) String() string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "ConsistencyLevel: %s\n", c.ConsistencyLevel)
	fmt.Fprintf(builder, "GracePeriod: %s\n", c.GracePeriod)
	fmt.Fprintf(builder, "Timeout: %s\n", c.Timeout)
	fmt.Fprintf(builder, "RF: %d\n", c.RF)

	fmt.Fprintln(builder, "Ring:")
	for id, node := range c.Ring {
		fmt.Fprintf(builder, "  ID %d -> %s\n", id, node)
	}

	return builder.String()
}

type ConfigManager struct {
	filepath string
	Data     Config
}

func newConfigManager(path string) *ConfigManager {
	if !filepath.IsAbs(path) {
		panic(errors.New("the provided path is not an absolute path"))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var data Config
	json.Unmarshal(file, &data)
	return &ConfigManager{filepath: path, Data: data}
}

func (configManager *ConfigManager) findNodeByPort(port int) (int, *Node) {
	for id, node := range configManager.Data.Ring {
		if node.Port == port {
			return id, &node
		}
	}
	panic(fmt.Errorf("could not find node with port %d", port))
}

func (configManager *ConfigManager) findNodeById(id int) *Node {
	node, exists := configManager.Data.Ring[id]
	if !exists {
		panic(fmt.Errorf("could not find node with id %d", id))
	}
	return &node
}
