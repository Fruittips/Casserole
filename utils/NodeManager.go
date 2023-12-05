package utils

import (
	"casserole/utils/cht"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

const CONFIG_PATH = "./config.json"
const DB_FILE_PATH = "./dbFiles/node-%d.json"
const HH_FILE_PATH = "./hintedHandoffs/node-%d.json"

type NodeId string

type Node struct {
	Id     NodeId
	Port   int
	isDead bool
}

func (n *Node) IsDead() bool {
	return n.isDead
}

func (n *Node) MakeDead() {
	n.isDead = true
}

func (n *Node) MakeAlive() {
	n.isDead = false
}

func (n Node) String() string {
	fmtStr := fmt.Sprintf("N%v: Port %d, Status: ", n.Id, n.Port)
	if n.IsDead() {
		return fmtStr + "DEAD"
	}
	return fmtStr + "LIVE"
}

// Manager for local node: Keeps track of status of other nodes
type NodeManager struct {
	LocalId              NodeId
	Quorum               int
	Nodes                map[NodeId](*Node)
	DatabaseManager      *DatabaseManager
	HintedHandoffManager *HintedHandoffManager
	//ReadRepairsManager   *ReadRepairsManager
	cht       *cht.CHashTable
	sysConfig *sysConfig
}

func NewNodeManager(port int, isSingleNode bool) *NodeManager {
	// Load filepaths
	configPath, err := filepath.Abs(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Error initialising NodeManager: %v", err)
	}
	dbPath, err := filepath.Abs(fmt.Sprintf(DB_FILE_PATH, port))
	if err != nil {
		log.Fatalf("Error initialising DatabaseManager: %v", err)
	}
	hhPath, err := filepath.Abs(fmt.Sprintf(HH_FILE_PATH, port))
	if err != nil {
		log.Fatalf("Error initialising HintedHandoffManager: %v", err)
	}

	// Identify ID of this node, generate nodeList for cht, and generate node map
	sysConfig, err := loadConfig(configPath, isSingleNode)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	myId := ""
	nodeLs := make([]cht.NodeId, 0)
	nodeMap := make(map[NodeId]*Node, 0)
	for id, nodeData := range sysConfig.Nodes {
		identifier := cht.NodeId(fmt.Sprintf("%d", nodeData.Port))
		nodeLs = append(nodeLs, identifier)
		nodeMap[NodeId(id)] = &Node{
			Id:     NodeId(id),
			Port:   nodeData.Port,
			isDead: nodeData.IsDead,
		}

		if nodeData.Port == port {
			if myId != "" {
				log.Fatalf("IDs %v, %v have the same port %d", id, myId, port)
			}
			myId = id
		}
	}

	if myId == "" {
		log.Fatalf("Could not find port %d in config.", port)
	}

	// Load Database Manager
	dbMgr, err := newDatabaseManager(dbPath)
	if err != nil {
		log.Fatalf("Error loading DatabaseManager: %v", err)
	}

	// Load HintedHandoffManager
	hhMgr, err := newHintedHandoffManager(hhPath)
	if err != nil {
		log.Fatalf("Error loading HintedHandoffManager: %v", err)
	}

	var quorum int

	switch sysConfig.ConsistencyLevel {
	case "ONE":
		quorum = 1
		break
	case "TWO":
		quorum = 2
		break
	case "THREE":
		quorum = 3
		break
	case "QUORUM":
		quorum = sysConfig.RF/2 + 1
		break
	case "ALL":
		quorum = sysConfig.RF
		break
	default:
		panic("Invalid consistency level")
	}

	if sysConfig.RF > len(sysConfig.Nodes) {
		panic("RF must be less than or equal to the number of nodes")
	}
	if sysConfig.RF < quorum {
		panic("RF must be greater than or equal to the quorum")
	}
	if len(sysConfig.Nodes) < quorum {
		panic("Number of nodes must be greater than or equal to the quorum")
	}

	return &NodeManager{
		LocalId:              NodeId(myId),
		Quorum:               quorum,
		Nodes:                nodeMap,
		DatabaseManager:      dbMgr,
		HintedHandoffManager: hhMgr,
		cht:                  cht.NewCHashTable(nodeLs),
		sysConfig:            sysConfig,
	}
}

// Returns a read-only reference of the configuration of this node manager
func (nm *NodeManager) GetConfig() sysConfig {
	return *nm.sysConfig
}

// Returns the local node
func (nm *NodeManager) Me() *Node {
	for id, node := range nm.Nodes {
		if id == nm.LocalId {
			return node
		}
	}
	panic("Local node not found in NodeManager.Nodes: Error in initialisation?")
}

// Returns the node associated with this port
func (nm *NodeManager) GetNodeByPort(port int) (*Node, error) {
	for _, node := range nm.Nodes {
		if node.Port == port {
			return node, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No node with port %d", port))
}

// Returns the node associated with this id
func (nm *NodeManager) GetNodeById(id NodeId) (*Node, error) {
	node, exists := nm.Nodes[id]
	if !exists {
		return nil, errors.New(fmt.Sprintf("No node with ID %v", id))
	}
	return node, nil
}

// Returns the node ports responsible for this key. Uses the system's config for RF.
func (nm *NodeManager) GetNodePortsForKey(key string) []int {
	portStrs := nm.cht.GetNodes(key, nm.sysConfig.RF)

	// Convert into ints
	ports := make([]int, 0)
	for _, portStr := range portStrs {
		port, err := strconv.Atoi(string(portStr))
		if err != nil {
			panic(fmt.Sprintf("Error converting %v to int: %v", portStr, err))
		}
		ports = append(ports, port)
	}
	return ports
}

// Returns the nodes responsible for this key. Uses the system's config for RF.
func (nm *NodeManager) GetNodesForKey(key string) [](*Node) {
	ports := nm.GetNodePortsForKey(key)
	nodes := make([](*Node), 0)
	for _, port := range ports {
		node, err := nm.GetNodeByPort(port)
		if err != nil {
			panic(fmt.Sprintf("Error getting port of node ID %v: %v", node.Id, err))
		}
		nodes = append(nodes, node)
	}
	return nodes
}
