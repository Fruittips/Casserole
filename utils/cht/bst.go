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
