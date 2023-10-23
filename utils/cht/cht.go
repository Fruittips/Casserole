/**
  `cht.CHashTable`: A local consistent hash table. Note that this only stores the hash IDs and corresponding node identifiers for each known node, and doesn't actually store any actual data.
  - `cht.NewCHashTable(nodeList []cht.NodeId) *CHashTable`: Creates a new hash table, given a list of cht.NodeId node identifiers (IP addresses, ports, etc). NodeId is an alias for a string.
  - `cht.AddNode(nodeId cht.NodeId)`: Adds a new node to the table.
  - `cht.RemoveNode(nodeId cht.NodeId)`: Deletes a node from the table.
  - `cht.GetNode(key string) NodeId`: Returns the node identifier for the node responsible for the given key.
*/

package cht

import (
	"fmt"
	"sync"
	"errors"

	"github.com/spaolacci/murmur3"
)

// Default seed for Murmur3 Hashing
const defaultSeed = 69

type HashId uint32 // Hash identifier for a node
type NodeId string // Identifier for a node (IP address, port, etc)

// A consistent hash table
type CHashTable struct {
	nodes *bst
	seed uint32
	rwlock *sync.RWMutex
}

// Given a list of cht.NodeId values, returns a pointer to a consistent hash table.
func NewCHashTable(nodeList []NodeId) *CHashTable {
	table := CHashTable{newBST(), uint32(defaultSeed), &sync.RWMutex{}}
	table.rwlock.Lock()
	for _, nodeId := range(nodeList) {
		table.nodes.Insert(
			getHashId(string(nodeId), table.seed),
			nodeId,
		)
	}
	table.rwlock.Unlock()
	
	return &table
}

// Adds a new node to the hash table. Will return an error if a duplicate node is added.
func (table *CHashTable) AddNode(nodeId NodeId) error {
	table.rwlock.Lock(); defer table.rwlock.Unlock()

	// Check for duplicates
	hashId := getHashId(string(nodeId), table.seed)
	if table.nodes.Search(hashId) == nodeId {
		return errors.New(fmt.Sprintf("cht.AddNode: Attempted to add duplicate node with NodeId %v", nodeId))
	}
	
	table.nodes.Insert(hashId, nodeId)
	return nil
}

// Removes a node from the hash table.
func (table *CHashTable) RemoveNode(nodeId NodeId) {
	table.rwlock.Lock()
	hashId := getHashId(string(nodeId), table.seed)
	table.nodes.Delete(hashId)
	table.rwlock.Unlock()
}

// Returns the ID of the node responsible for this key.
func (table *CHashTable) GetNode(key string) NodeId {
	table.rwlock.RLock()
	defer table.rwlock.RUnlock()
	dataId := getHashId(key, table.seed)
	return table.nodes.ClosestSmallerNode(dataId)
}


func getHashId(key string, seed uint32) HashId {
	return HashId(murmur3.Sum32WithSeed([]byte(key), seed) % 360)
}
