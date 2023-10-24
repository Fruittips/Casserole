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

type NodeId int

type Node struct {
	Id NodeId
	Port int
	isDead bool
}

func (n Node) IsDead() bool {
	// TODO: replace with proper way -- why can't we get rid of deadnode and just use ctrl-c?
	return n.isDead
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
	LocalId NodeId
	Quorum int
	Nodes map[NodeId]Node
	cht *cht.CHashTable
	sysConfig *sysConfig
}

func NewNodeManager(port int) *NodeManager {
	configPath, err := filepath.Abs(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Error initialising NodeManager: %v", err)
	}

	// Identify ID of this node, generate nodeList for cht, and generate node map
	sysConfig, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	myId := -1
	nodeLs := make([]cht.NodeId, 0)
	nodeMap := make(map[NodeId]Node, 0)
	for id, nodeData := range(sysConfig.Nodes) {
		identifier := cht.NodeId(fmt.Sprintf("%d", nodeData.Port))
		nodeLs = append(nodeLs, identifier)
		nodeMap[NodeId(id)] = Node{
			Id: NodeId(id),
			Port: nodeData.Port,
			isDead: nodeData.IsDead,
		}
		
		if nodeData.Port == port {
			if myId != -1 {
				log.Fatalf("IDs %v, %v have the same port %d", id, myId, port)
			}
			myId = id
		}
	}

	if myId == -1 {
		log.Fatalf("Could not find port %d in config.", port)
	}

	return &NodeManager{
		LocalId: NodeId(myId),
		Quorum: sysConfig.RF/2 + 1,
		Nodes: nodeMap,
		cht: cht.NewCHashTable(nodeLs),
		sysConfig: sysConfig,
	}
}

// Returns a read-only reference of the configuration of this node manager
func (nm *NodeManager) GetConfig() sysConfig {
	return *nm.sysConfig
}

// Returns the node associated with this port
func (nm *NodeManager) GetNodeByPort(port int) (*Node, error) {
	for _, node := range(nm.Nodes) {
		if node.Port == port {
			return &node, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No node with port %d", port))
}

// Returns the node associated with this id
func (nm *NodeManager) GetNodeById(id NodeId) (*Node, error) {
	node, exists := nm.Nodes[id]
	if !exists {
		return nil, errors.New(fmt.Sprintf("No node with ID %d", id))
	}
	return &node, nil
}

// Returns the node ports responsible for this key. Uses the system's config for RF.
func (nm *NodeManager) GetNodePortsForKey(key string) []int {
	portStrs := nm.cht.GetNodes(key, nm.sysConfig.RF)

	// Convert into ints
	ports := make([]int, 0)
	for _, portStr := range(portStrs) {
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
	for _, port := range(ports) {
		node, err := nm.GetNodeByPort(port)
		if err != nil {
			panic(fmt.Sprintf("Error getting port of node ID %d: %v", node.Id, err))
		}
		nodes = append(nodes, node)
	}
	return nodes
}


