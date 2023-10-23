package cht

import (
	"testing"
	"fmt"
)

// Get the list of HashIds for testing
func (tree *bst) hashIds() []HashId {
	return tree.hashIdsRec(tree.root)
}

func assertHashIdsMatches(hashIds []HashId, ints []int) bool {
	if len(hashIds) != len(ints) {
		return false
	}

	for i := range(hashIds) {
		if int(hashIds[i]) != ints[i] {
			return false
		}
	}
	return true
}

func (tree *bst) hashIdsRec(node *bstNode) []HashId {
	ls := make([]HashId, 0)
	if node == nil {
		return ls
	}
	if node.left != nil {
		ls = append(ls, tree.hashIdsRec(node.left)...)
	}
	
	ls = append(ls, node.key)
	if node.right != nil {
		ls = append(ls, tree.hashIdsRec(node.right)...)
	}
	return ls
}

func TestSearchAndInsert(t *testing.T) {
	tree := newBST()
	testInts := []int{0, 1, 2, 3, 4, 5}

	// Inserts
	for _, testInt := range(testInts) {
		tree.Insert(HashId(testInt), NodeId(fmt.Sprintf("value%d", testInt)))
	}
	
	// Test Search
	for _, testInt := range(testInts) {
		nodeId := tree.Search(HashId(testInt))
		expectedNodeId := NodeId(fmt.Sprintf("value%d", testInt))
		if nodeId != expectedNodeId {
			t.Fatalf("tree.Search(%d) = %v, expected %v.", testInt, nodeId, expectedNodeId)
		}
	}

	// Test Inorder Sorted-ness
	for i, hashId := range(tree.hashIds()) {
		if hashId != HashId(testInts[i]) {
			t.Fatalf("tree is not sorted. hashId = %d, testInts[i] = %d.", hashId, testInts[i])
		}
	}

	t.Log(tree.hashIds())
}

func TestSearchAndOutOfOrderInsert(t *testing.T) {
	tree := newBST()
	testInts := []int{0, 1, 2, 3, 4, 5}

	// Inserts
	for i := len(testInts)-1; i >= 0; i-- {
		testInt := testInts[i]
		tree.Insert(HashId(testInt), NodeId(fmt.Sprintf("value%d", testInt)))
	}
	
	// Test Search
	for _, testInt := range(testInts) {
		nodeId := tree.Search(HashId(testInt))
		expectedNodeId := NodeId(fmt.Sprintf("value%d", testInt))
		if nodeId != expectedNodeId {
			t.Fatalf("tree.Search(%d) = %v, expected %v.", testInt, nodeId, expectedNodeId)
		}
	}

	// Test Inorder Sorted-ness
	for i, hashId := range(tree.hashIds()) {
		if hashId != HashId(testInts[i]) {
			t.Fatalf("tree is not sorted. hashId = %d, testInts[i] = %d.", hashId, testInts[i])
		}
	}

	t.Log(tree.hashIds())

	
	tree = newBST()
	testInts = []int{1,2,3,5,6,33,42,99}
	testIntsRandomOrder := []int{6,42,3,1,2,5,99,33}

	// Inserts
	for _, testInt := range(testIntsRandomOrder) {
		tree.Insert(HashId(testInt), NodeId(fmt.Sprintf("value%d", testInt)))
	}
	
	// Test Search
	for _, testInt := range(testInts) {
		nodeId := tree.Search(HashId(testInt))
		expectedNodeId := NodeId(fmt.Sprintf("value%d", testInt))
		if nodeId != expectedNodeId {
			t.Fatalf("tree.Search(%d) = %v, expected %v.", testInt, nodeId, expectedNodeId)
		}
	}

	// Test Inorder Sorted-ness
	for i, hashId := range(tree.hashIds()) {
		if hashId != HashId(testInts[i]) {
			t.Fatalf("tree is not sorted. hashId = %d, testInts[i] = %d.", hashId, testInts[i])
		}
	}

	t.Log(tree.hashIds())
}


func TestDeletes(t *testing.T) {
	tree := newBST()
	testInts := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	// Inserts
	for _, testInt := range(testInts) {
		tree.Insert(HashId(testInt), NodeId(fmt.Sprintf("value%d", testInt)))
	}
	
	// Test Search
	for _, testInt := range(testInts) {
		nodeId := tree.Search(HashId(testInt))
		expectedNodeId := NodeId(fmt.Sprintf("value%d", testInt))
		if nodeId != expectedNodeId {
			t.Fatalf("tree.Search(%d) = %v, expected %v.", testInt, nodeId, expectedNodeId)
		}
	}
	if !assertHashIdsMatches(tree.hashIds(), testInts) {
		t.Log(tree.hashIds())
		t.Fatalf("Inserts failed.")
	}

	// Test Deletes
	tree.Delete(HashId(11))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,2,3,4,5,6,7,8,9,10,12}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(11) failed.")
	}
	
	tree.Delete(HashId(8))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,2,3,4,5,6,7,9,10,12}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(8) failed.")
	}
	
	tree.Delete(HashId(10))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,2,3,4,5,6,7,9,12}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(10) failed.")
	}
	
	tree.Delete(HashId(9))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,2,3,4,5,6,7,12}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(9) failed.")
	}

	tree.Delete(HashId(12))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,2,3,4,5,6,7}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(12) failed.")
	}
	
	tree.Delete(HashId(4))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,2,3,5,6,7}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(4) failed.")
	}
	
	tree.Delete(HashId(2))
	if !assertHashIdsMatches(tree.hashIds(), []int{0,1,3,5,6,7}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(2) failed.")
	}
	
	tree.Delete(HashId(0))
	if !assertHashIdsMatches(tree.hashIds(), []int{1,3,5,6,7}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(0) failed.")
	}
	
	tree.Delete(HashId(1))
	if !assertHashIdsMatches(tree.hashIds(), []int{3,5,6,7}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(1) failed.")
	}
	
	tree.Delete(HashId(3))
	if !assertHashIdsMatches(tree.hashIds(), []int{5,6,7}) {
		t.Log(tree.hashIds())
		t.Fatalf("tree.Delete(3) failed.")
	}
}

func TestClosestSmallerNode(t *testing.T) {
	tree := newBST()
	testInts := []int{0, 2, 4, 6, 8, 10, 12}

	// Insert Nodes
	for _, testInt := range(testInts) {
		tree.Insert(HashId(testInt), NodeId(fmt.Sprintf("value%d", testInt)))
	}
	if !assertHashIdsMatches(tree.hashIds(), testInts) {
		t.Log(tree.hashIds())
		t.Fatalf("Inserts failed.")
	}

	// Test Equal Data Keys
	for _, testInt := range(testInts) {
		key := HashId(testInt)
		closestNodeId := tree.ClosestSmallerNode(key)
		expectedNodeId := NodeId(fmt.Sprintf("value%d", testInt))
		if closestNodeId != expectedNodeId {
			t.Fatalf("Closest NodeId for HashId %v was: %v; expected: %v", key, closestNodeId, expectedNodeId)
		}
	}

	// Test Nonequal Data Keys
	keys             := []HashId{1, 3, 5, 7, 9, 11, 13, 15, 17}
	expectedNodeHashIds := []int{0, 2, 4, 6, 8, 10, 12, 12, 12}
	expectedNodeIds := make([]NodeId, 0)
	for _, nodeHashId := range(expectedNodeHashIds) {
		expectedNodeIds = append(expectedNodeIds, NodeId(fmt.Sprintf("value%d", nodeHashId)))
	}

	for i, key := range(keys) {
		closestNodeId := tree.ClosestSmallerNode(key)
		expectedNodeId := expectedNodeIds[i]
		if closestNodeId != expectedNodeId {
			t.Fatalf("Closest NodeId for HashId %v was: %v; expected: %v", key, closestNodeId, expectedNodeId)
		}
	}

	// Test if data keys lower than minimum node HashId get assigned circularly
	tree = newBST()
	nodeHashIds := []HashId{5, 10, 15, 20, 25, 30, 50}
	for _, hashId := range(nodeHashIds) {
		tree.Insert(hashId, NodeId(fmt.Sprintf("value%d", hashId)))
	}

	keys              = []HashId{1, 2, 5, 6, 9, 10, 22, 24, 33, 39, 60}
	expectedNodeHashIds = []int{50,50, 5, 5, 5, 10, 20, 20, 30, 30, 50}
	expectedNodeIds = make([]NodeId, 0)
	for _, nodeHashId := range(expectedNodeHashIds) {
		expectedNodeIds = append(expectedNodeIds, NodeId(fmt.Sprintf("value%d", nodeHashId)))
	}

	for i, key := range(keys) {
		closestNodeId := tree.ClosestSmallerNode(key)
		expectedNodeId := expectedNodeIds[i]
		if closestNodeId != expectedNodeId {
			t.Fatalf("Closest NodeId for HashId %v was: %v; expected: %v", key, closestNodeId, expectedNodeId)
		}
	}
	
}
