package searchtree

import (
	"fmt"
	"sync"
)

type Item string

type Node struct {
	Key   string
	Value Item
	Left  *Node
	Right *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("Node{Key=%s,Value=%s}", n.Key, n.Value)
}

type BinarySearchTree struct {
	lock sync.RWMutex
	Root *Node
}

func (bst *BinarySearchTree) String() string {
	return fmt.Sprintf("BST{Root=%s}", bst.Root)
}

func (bst *BinarySearchTree) Insert(key string, value Item) {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	n := &Node{key, value, nil, nil}

	if bst.Root == nil {
		bst.Root = n
	} else {
		insertNode(bst.Root, n)
	}
}

func insertNode(node, newNode *Node) {
	if newNode.Key < node.Key {
		if node.Left == nil {
			node.Left = newNode
		} else {
			insertNode(node.Left, newNode)
		}
	} else {
		if node.Right == nil {
			node.Right = newNode
		} else {
			insertNode(node.Right, newNode)
		}
	}
}

func (bst *BinarySearchTree) InOrderTraverse(f func(Item)) {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	inOrderTraverse(bst.Root, f)
}

func inOrderTraverse(n *Node, f func(Item)) {
	if n != nil {
		inOrderTraverse(n.Left, f)
		f(n.Value)
		inOrderTraverse(n.Right, f)
	}
}

// PreOrderTraverse visits all nodes with pre-order traversing
func (bst *BinarySearchTree) PreOrderTraverse(f func(Item)) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	preOrderTraverse(bst.Root, f)
}

// internal recursive function to traverse pre order
func preOrderTraverse(n *Node, f func(Item)) {
	if n != nil {
		f(n.Value)
		preOrderTraverse(n.Left, f)
		preOrderTraverse(n.Right, f)
	}
}

// PostOrderTraverse visits all nodes with post-order traversing
func (bst *BinarySearchTree) PostOrderTraverse(f func(Item)) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	postOrderTraverse(bst.Root, f)
}

// internal recursive function to traverse post order
func postOrderTraverse(n *Node, f func(Item)) {
	if n != nil {
		postOrderTraverse(n.Left, f)
		postOrderTraverse(n.Right, f)
		f(n.Value)
	}
}

// Min returns the Item with min Value stored in the tree
func (bst *BinarySearchTree) Min() *Item {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	n := bst.Root
	if n == nil {
		return nil
	}
	for {
		if n.Left == nil {
			return &n.Value
		}
		n = n.Left
	}
}

// Max returns the Item with max Value stored in the tree
func (bst *BinarySearchTree) Max() *Item {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	n := bst.Root
	if n == nil {
		return nil
	}
	for {
		if n.Right == nil {
			return &n.Value
		}
		n = n.Right
	}
}

// Search returns true if the Item t exists in the tree
func (bst *BinarySearchTree) Search(key string) bool {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	return search(bst.Root, key)
}

// internal recursive function to search an item in the tree
func search(n *Node, key string) bool {
	if n == nil {
		return false
	}
	if key < n.Key {
		return search(n.Left, key)
	}
	if key > n.Key {
		return search(n.Right, key)
	}
	return true
}

// Remove removes the Item with key `key` from the tree
func (bst *BinarySearchTree) Remove(key string) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	remove(bst.Root, key)
}

// internal recursive function to remove an item
func remove(node *Node, key string) *Node {
	if node == nil {
		return nil
	}
	if key < node.Key {
		node.Left = remove(node.Left, key)
		return node
	}
	if key > node.Key {
		node.Right = remove(node.Right, key)
		return node
	}
	// key == node.key
	if node.Left == nil && node.Right == nil {
		node = nil
		return nil
	}
	if node.Left == nil {
		node = node.Right
		return node
	}
	if node.Right == nil {
		node = node.Left
		return node
	}
	LeftmostRightside := node.Right
	for {
		//find smallest Value on the Right side
		if LeftmostRightside != nil && LeftmostRightside.Left != nil {
			LeftmostRightside = LeftmostRightside.Left
		} else {
			break
		}
	}
	node.Key, node.Value = LeftmostRightside.Key, LeftmostRightside.Value
	node.Right = remove(node.Right, node.Key)
	return node
}

// Search returns the item if it exists, else returns an empty string
func (bst *BinarySearchTree) Get(key string) Item {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	return get(bst.Root, key)
}

// internal recursive function to search an item in the tree
func get(n *Node, key string) Item {
	if n == nil {
		return ""
	}
	if key < n.Key {
		return get(n.Left, key)
	}
	if key > n.Key {
		return get(n.Right, key)
	}
	return n.Value
}
