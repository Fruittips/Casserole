package utils

import (
	"casserole/utils/cht"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
)

// A local record of a node
type Node struct {
	Id int
	Port   int
	IsDead bool
}

func (n Node) String() string {
	deadStatus := "Alive"
	if n.IsDead {
		deadStatus = "Dead"
	}
	return fmt.Sprintf("Port: %d, Status: %s", n.Port, deadStatus)
}

// Manager for the local node.
// This keeps track of the status of all other nodes, manages the database and hinted handoffs.
type NodeManager struct {
	Id int
	Quorum int
	Ring map[int]Node
	ConfigManager *ConfigManager // Default configuration
	DatabaseManager *DatabaseManager
	HintedHandoffManager *HintedHandoffManager
	cht cht.CHashTable
}

// Initialises the node, loading all local data 'owned' by this node
func NewNodeManager(port int) *NodeManager {
	relativePath := "./config.json"
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	// Identify the ID of this port, based on the configuration
	configManager := newConfigManager(absolutePath)
	myId := -1

	// Identify ID and generate nodeList for cht
	nodeList := make([]cht.NodeId, 0)
	for id, node := range configManager.defaultNodes {
		nodeIdentifier := cht.NodeId(fmt.Sprintf("%d", node.Port))
		nodeList = append(nodeList, nodeIdentifier)
		if node.Port == port {
			myId = id
		}
	}
	if myId == -1 {
		panic(fmt.Sprintf("Could not find port %d in config", port))
	}
	

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
		Quorum:               configManager.RF/2 + 1,
		Ring: configManager.defaultNodes,
		ConfigManager:        configManager,
		DatabaseManager:      databaseManager,
		HintedHandoffManager: hintedHandoffManager,
		cht: *cht.NewCHashTable(nodeList),
	}
}

// Returns the node associated with port
func (nMgr *NodeManager) GetNodeByPort(port int) (*Node, error) {
	for _, node := range(nMgr.Ring) {
		if node.Port == port {
			return &node, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No node with port %d", port))
}

// Returns the node associated with this ID in the configuration
func (nMgr *NodeManager) GetNodeById(id int) (*Node, error) {
	node, exists := nMgr.Ring[id]
	if !exists {
		return nil, errors.New(fmt.Sprintf("No node with ID %d", id))
	}
	return &node, nil
}

// Returns the nodes responsible for this key. This uses the replication factor in the configuration
func (nMgr *NodeManager) GetNodesForKey(key string) [](*Node) {
	nodeIdentifiers := nMgr.cht.GetNodes(key, nMgr.ConfigManager.RF)
	nodes := make([](*Node), 0)
	for _, identifier := range(nodeIdentifiers) {
		port, err := strconv.Atoi(string(identifier))
		if err != nil {
			panic(fmt.Sprintf("Error converting %v to int: %v", identifier, err))
		}
		node, err:= nMgr.GetNodeByPort(port)
		if err != nil {
			panic(fmt.Sprintf("Error getting nodes by port %d: %v", port, node))
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// Returns the node identifiers of the nodes responsible for this key. This uses the RF in the config
func (nMgr *NodeManager) GetNodePortsForKey(key string) [](int) {
	nodeIdentifiers := nMgr.cht.GetNodes(key, nMgr.ConfigManager.RF)
	ports := make([]int, 0)
	for _, identifier := range(nodeIdentifiers) {
		port, err := strconv.Atoi(string(identifier))
		if err != nil {
			panic(fmt.Sprintf("Error converting %v to int: %v", identifier, err))
		}
		ports = append(ports, port)
	}
	return ports
}

func (nm *NodeManager) Me() Node {
	me, err := nm.GetNodeById(nm.Id)
	if err != nil {
		panic(err)
	}
	return *me
}

