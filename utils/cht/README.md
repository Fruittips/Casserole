# cht

`cht.CHashTable` is a local consistent hash table that stores the hash IDs and corresponding node identifiers for each known node.
- The "node identifier" could be an IP address, a port or anything that uniquely identifies the node.
- Note that this doesn't actually store any data, the `CHashTable` only knows *which* node ID a given data key belongs to.

The consistent hash table uses a binary search tree `bst` in the background, and this is also used to identify the node to which data keys should be alloted to.

## Interface
### `cht.NewCHashTable(nodeList []cht.NodeId) *cht.CHashTable`
Creates a new hash table.
- `nodeList`: A slice of `cht.NodeId` node identifiers.
  - Note that `cht.NodeId` is just an alias for a string, an example node ID would be `cht.NodeId("192.168.0.1")` or `cht.NodeId("6969")`.
- This returns a pointer to a `CHashTable` instance.

### `table.GetNode(key string) cht.NodeId`
Returns the `cht.NodeId` node identifier for a given `key` string.

### `table.AddNode(nodeId cht.NodeId) error`
Adds a node ID to the table. Returns an error if the given `nodeId` is already in the table.

### `table.RemoveNode(nodeId cht.NodeId)`
Remove `nodeId` from the table.

## Example Usage
```go
import "cht"

func main() {
	// A nodelist for processes running on the following ports: 6969, 8080, 80 and 443.
	nodeList := []cht.NodeId{
		"6969",
		"8080",
		"80",
		"443",
	}
	
	// Initialise the table
	table := cht.NewCHashTable(nodeList)
	
	// Data tuple (key, value)
	key := "02.069"
	value := "Duck-flavoured Potatoes and their Repercussions on World War XI"
	
	// Identify which node is responsible for this data tuple
	nodeId := table.GetNode(key)
}
```
- Note that the actual node ID returned could vary, based on the `defaultSeed` value set in `cht.go`.
