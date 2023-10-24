package cht

import (
	"testing"
)

func TestNodes(t *testing.T) {
	table := NewCHashTable(make([]NodeId, 0))
	err := table.AddNode(NodeId("6969"))
	if err != nil {
		t.Fatalf("Error adding node 6969: %v", err)
	}
	err = table.AddNode(NodeId("8080"))
	if err != nil {
		t.Fatalf("Error adding node 8080: %v", err)
	}
	err = table.AddNode(NodeId("8081"))
	if err != nil {
		t.Fatalf("Error adding node 8081: %v", err)
	}
	err = table.AddNode(NodeId("8079"))
	if err != nil {
		t.Fatalf("Error adding node 8079: %v", err)
	}
	err = table.AddNode(NodeId("8080"))
	if err == nil {
		t.Fatalf("Error: Expected error adding duplicate node 8080, but no error occurred.")
	}
	err = table.AddNode(NodeId("6969"))
	if err == nil {
		t.Fatalf("Error: Expected error adding duplicate node 6969, but no error occurred.")
	}	
}

func TestGetNode(t *testing.T) {
	// Note: Most of the necessary tests here are already done in the binary search tree tests
	// For instance, the value of the closest node based on the data key is already tested there.
	// Here, we merely test if GetNode works as we expect.
	
	table := NewCHashTable(make([]NodeId, 0))
	err := table.AddNode(NodeId("6969"))
	if err != nil {
		t.Fatalf("Error adding node 6969: %v", err)
	}
	err = table.AddNode(NodeId("8080"))
	if err != nil {
		t.Fatalf("Error adding node 8080: %v", err)
	}
	err = table.AddNode(NodeId("8081"))
	if err != nil {
		t.Fatalf("Error adding node 8081: %v", err)
	}
	err = table.AddNode(NodeId("8079"))
	if err != nil {
		t.Fatalf("Error adding node 8079: %v", err)
	}

	// Here, a data key "6969" will be hashed to the exact same value anyway, so we use this to test.
	nodeId := table.GetNode("6969")
	if nodeId != NodeId("6969") {
		t.Fatalf("NodeIDs don't match for value \"6969\". answer: %v, expected %v", nodeId, NodeId("6969"))
	}

	//TODO: add tests based on the hashing
	
}

func TestGetNodes(t *testing.T) {
	table := NewCHashTable(make([]NodeId, 0))
	err := table.AddNode(NodeId("6969"))
	if err != nil {
		t.Fatalf("Error adding node 6969: %v", err)
	}
	err = table.AddNode(NodeId("8080"))
	if err != nil {
		t.Fatalf("Error adding node 8080: %v", err)
	}
	err = table.AddNode(NodeId("8081"))
	if err != nil {
		t.Fatalf("Error adding node 8081: %v", err)
	}
	err = table.AddNode(NodeId("8079"))
	if err != nil {
		t.Fatalf("Error adding node 8079: %v", err)
	}

	// Test GetNodes by checking the size of the list: The main checks are done in bst_test
	nodeIds := table.GetNodes("asdf", 0)
	if len(nodeIds) != 0 {
		t.Fatalf("Expected 0 nodeIds, instead got %v", len(nodeIds))
	}
	
	nodeIds = table.GetNodes("asdf", 1)
	if len(nodeIds) != 1 {
		t.Fatalf("Expected 1 nodeId, instead got %v", len(nodeIds))
	}
	
	nodeIds = table.GetNodes("asdf", 2)
	if len(nodeIds) != 2 {
		t.Fatalf("Expected 2 nodeIds, instead got %v", len(nodeIds))
	}
	
	nodeIds = table.GetNodes("asdf", 3)
	if len(nodeIds) != 3 {
		t.Fatalf("Expected 3 nodeIds, instead got %v", len(nodeIds))
	}

	// (All nodes)
	nodeIds = table.GetNodes("asdf", 4)
	if len(nodeIds) != 4 {
		t.Fatalf("Expected 4 nodeIds, instead got %v", len(nodeIds))
	}

	// Should reset to 4
	nodeIds = table.GetNodes("asdf", 5)
	if len(nodeIds) != 4 {
		t.Fatalf("Expected 3 nodeIds, instead got %v", len(nodeIds))
	}
	//TODO: add tests based on the hashing
}
