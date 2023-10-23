package utils

import (
	"fmt"
	"path/filepath"
)

type Node struct {
	Port   int  `json:"port"` //e.g. "http://localhost:3000"
	IsDead bool `json:"isDead"`
}

type NodeManager struct {
	id                   int
	ConfigManager        *ConfigManager
	DatabaseManager      *DatabaseManager
	HintedHandoffManager *HintedHandoffManager
}

func (n Node) String() string {
	deadStatus := "Alive"
	if n.IsDead {
		deadStatus = "Dead"
	}
	return fmt.Sprintf("Port: %d, Status: %s", n.Port, deadStatus)
}

func NewNodeManager(port int) *NodeManager {
	relativePath := "./config.json"
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	configManager := newConfigManager(absolutePath)
	myId, _ := configManager.findNodeByPort(port)

	databaseFilepath := fmt.Sprintf("./dbFiles/node-%d.json", myId)
	absolutePathDb, err := filepath.Abs(databaseFilepath)
	if err != nil {
		panic(err)
	}

	hintedHandoffFilepath := fmt.Sprintf("./hintedHandoffs/node-%d.json", myId)
	absolutePathHh, err := filepath.Abs(hintedHandoffFilepath)
	if err != nil {
		panic(err)
	}

	databaseManager := newDatabaseManager(absolutePathDb)
	hintedHandoffManager := newHintedHandoffManager(absolutePathHh)

	return &NodeManager{
		id:                   myId,
		ConfigManager:        configManager,
		DatabaseManager:      databaseManager,
		HintedHandoffManager: hintedHandoffManager,
	}
}

func (nm *NodeManager) Me() Node {
	me := nm.ConfigManager.findNodeById(nm.id)
	return *me
}
