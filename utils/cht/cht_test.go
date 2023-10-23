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
