package utils

import (
	"fmt"
	"path/filepath"
)

type RequestType int

const (
	Read  RequestType = iota // 0 for read request
	Write                    // 1 for write request
)

var RequestTypeStr = map[RequestType]string{
	Read:  "read",
	Write: "write"}

type Node struct {
	Port   int  `json:"port"` //e.g. "http://localhost:3000"
	IsDead bool `json:"isDead"`
}

type NodeManager struct {
	Id                   int
	Quorum               int
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
		Id:                   myId,
		Quorum:               configManager.Data.RF/2 + 1,
		ConfigManager:        configManager,
		DatabaseManager:      databaseManager,
		HintedHandoffManager: hintedHandoffManager,
	}
}

func (nm *NodeManager) Me() Node {
	me := nm.ConfigManager.FindNodeById(nm.Id)
	return *me
}
