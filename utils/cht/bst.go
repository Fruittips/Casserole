/**
  A binary search tree implementation which associates the key with a HashId
  and value with a NodeId.
*/

package cht

type bstNode struct {
	key HashId
	value NodeId
	parent *bstNode
	left   *bstNode
	right  *bstNode
}

type bst struct {
	root *bstNode
}

func newBST() *bst {
	return &bst{nil}
}

// Searches the tree for the given HashId. Returns an empty string if not found.
func (tree *bst) Search(key HashId) NodeId {
	node := tree.searchRec(tree.root, key)
	if node == nil {
		return ""
	} else {
		return node.value
	}
}

// Returns the NodeId of the node that has the **closest** HashId that is smaller
// than the given HashId.
// If no node has a smaller HashId, returns the node with the largest HashId
// since this is organised in a ring structure
// If there are no nodes in the hash table, returns an empty NodeId
func (tree *bst) ClosestSmallerNode(key HashId) NodeId {
	closestNode := tree.closestSmallerNodeRec(tree.root, key)
	if closestNode == nil {
		// Either tree root is nil, or we found no nodes smaller than key
		if tree.root == nil {
			return ""
		}
		// Get maximum
		closestNode = tree.root
		for closestNode.right != nil {
			closestNode = closestNode.right
		}
	}
	return closestNode.value
}

// Given a HashId of a node, return a list of HashIds of the successors of that node.
// If successorCount exceeds len(tree)-2, returns all available successors.
func (tree *bst) GetSuccessors(key HashId, successorCount int) []HashId {
	successors := make([]HashId, 0)

	if successorCount == 0 {
		return successors
	}

	if successorCount > tree.Len() - 2 {
		// Return all available successors
		successorCount = tree.Len() - 2
	}
	return successors
}

// Returns the successor of a given node
func (tree *bst) successor(node *bstNode) *bstNode {
	return node
}

// Insert key into tree
func (tree *bst) Insert(key HashId, value NodeId) {
	curNode := tree.root

	// Identify parent of key
	var parent *bstNode
	for curNode != nil {
		parent = curNode
		if key < curNode.key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
	}

	// Here, curNode is nil
	curNode = &bstNode{
		key,
		value,
		parent,
		nil, nil,
	}

	if parent == nil {
		tree.root = curNode
	} else {
		// identify whether curNode is left or right of parent
		if curNode.key < parent.key {
			parent.left = curNode
		} else {
			parent.right = curNode
		}
	}
}

// Deletes an item from the tree
func (tree *bst) Delete(key HashId) {
	node := tree.searchRec(tree.root, key)
	tree.deleteNode(node)
}

// Returns a slice of all HashIds in the tree.
func (tree *bst) HashIds() []HashId {
	return tree.hashIdsRec(tree.root)
}

// Returns the number of nodes in the tree
func (tree *bst) Len() int {
	return len(tree.HashIds())
}

// Recursively searches tree for key. Returns either the bstNode or nil.
func (tree *bst) searchRec(node *bstNode, key HashId) *bstNode {
	if node == nil || key == node.key {
		return node
	}

	if key < node.key {
		return tree.searchRec(node.left, key)
	} else {
		return tree.searchRec(node.right, key)
	}
}

func (tree *bst) closestSmallerNodeRec(node *bstNode, key HashId) *bstNode {
	// CASES:
	// key > curNode.key: Closest node is in the right subtree, or is curNode
	// - Return max(curNode, closestSmallerNodeRec(curNode.right))
	// - if no right subtree, return curNode
	// key < curNode.key: Closest node is in the left subtree, or is max node
	// - Return closestSmallerNodeRec(curNode.left)
	// - if no left subtree, return nil
	// key == curNode.key: This IS the closest node

	if node == nil {
		return nil
	}

	if key > node.key {
		closestRightChild := tree.closestSmallerNodeRec(node.right, key)
		if closestRightChild == nil {
			// Either node.right is nil, or there's no larger child smaller than key
			return node
		} else {
			return closestRightChild
		}
	} else if key < node.key {
		return tree.closestSmallerNodeRec(node.left, key)
	} else {
		return node
	}
}

func (tree *bst) deleteNode(node *bstNode) {
	if node == nil {
		return
	}

	// If node has no children, just remove the node
	if node.left == nil && node.right == nil {
		if node == node.parent.left {
			node.parent.left = nil
		} else {
			node.parent.right = nil
		}
		return
	}

	// If node has one child, we can elevate the child
	if node.left == nil {
		tree.transplant(node, node.right)
		return
	} else if node.right == nil {
		tree.transplant(node, node.left)
		return
	}

	// If node has two children, we need to find the successor to transplant
	// This successor is the leftmost element of the subtree rooted at node.right
	successor := node.right
	for successor.left != nil {
		successor = successor.left
	}

	if successor.parent == node {
		// If successor is node's right child, then replace
		tree.transplant(node, successor)
		successor.left = node.left
		successor.left.parent = successor
	} else {
		// Otherwise
		tree.transplant(successor, successor.right)
		successor.right = node.right
		successor.right.parent = successor
	}
}

// Replaces subtree with root u with subtree with root v
func (tree *bst) transplant(u, v *bstNode) {
	if u.parent == nil {
		// u is root node
		tree.root = v
	} else if u == u.parent.left {
		// u is a left child
		u.parent.left = v
	} else if u == u.parent.right {
		// u is a right child
		u.parent.right = v
	}

	if v != nil {
		// Update v to recognise parent of u
		v.parent = u.parent
	}
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
