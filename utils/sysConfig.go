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

// System-level configuration options
type sysConfig struct {
	ConsistencyLevel string                `json:"consistencyLevel"`
	GracePeriod      time.Duration         `json:"gracePeriod" `
	Timeout          time.Duration         `json:"timeout"`
	RF               int                   `json:"rf"`
	Nodes            map[string]nodeConfig `json:"ring"`
}

// Temporary struct to load from config
type nodeConfig struct {
	Port   int  `json:"port"`
	IsDead bool `json:"isDead"`
}

// Loads configuration from a given path
func loadConfig(path string, isSingleNode bool) (*sysConfig, error) {
	if !filepath.IsAbs(path) {
		return nil, errors.New(fmt.Sprintf("Expected absolute path, was given %v", path))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read file %v, error: %v", path, err))
	}

	var configData sysConfig
	err = json.Unmarshal(file, &configData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal JSON file %v, error: %v", path, err))
	}

	if isSingleNode {
		configData.RF = 1
		configData.Nodes = map[string]nodeConfig{
			"3000": {
				Port:   3000,
				IsDead: false,
			},
		}
	}

	return &configData, nil
}

func (conf *sysConfig) String() string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "ConsistencyLevel: %s\n", conf.ConsistencyLevel)
	fmt.Fprintf(builder, "GracePeriod: %s\n", conf.GracePeriod)
	fmt.Fprintf(builder, "Timeout: %s\n", conf.Timeout)
	fmt.Fprintf(builder, "RF: %d\n", conf.RF)

	fmt.Fprintln(builder, "Ring:")
	for id, node := range conf.Nodes {
		fmt.Fprintf(builder, "  ID %v -> %s\n", id, node.String())
	}

	return builder.String()
}

func (nodeConf *nodeConfig) String() string {
	if nodeConf.IsDead {
		return fmt.Sprintf("Port: %d, Status: DEAD", nodeConf.Port)
	}
	return fmt.Sprintf("Port: %d, Status: LIVE", nodeConf.Port)
}
