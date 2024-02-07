package avl

import "cmp"

type C interface {
	comparable
	cmp.Ordered
}

type AVLTree[T C] struct {
	nodeCount uint32
	root      *Node[T]
}

type Node[T C] struct {
	height        uint32
	key           T
	val           interface{}
	left          *Node[T]
	right         *Node[T]
	balanceFactor int
}

func (a *AVLTree[C]) Insert(key C, val interface{}) bool {
	if !a.Contains(key) {
		a.root = a.insert(a.root, key, val)
		a.nodeCount++
		return true
	}

	return false
}

func (a *AVLTree[C]) Remove(key C) bool {
	if a.Contains(key) {
		a.root = a.remove(a.root, key)
		a.nodeCount--
		return true
	}

	return false
}

func (a *AVLTree[C]) remove(n *Node[C], key C) *Node[C] {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = a.remove(n.left, key)
	} else if key > n.key {
		n.right = a.remove(n.right, key)
	} else {
		if n.left == nil {
			return n.right
		}

		if n.right == nil {
			return n.left
		}

		successorKey, successorVal := a.findMax(n.left)
		n.key = successorKey
		n.val = successorVal

		n.left = a.remove(n.left, successorKey)
	}

	a.update(n)

	return a.balance(n)
}

func (a *AVLTree[C]) findMax(n *Node[C]) (key C, val interface{}) {
	tmp := n
	for tmp.right != nil {
		tmp = tmp.right
	}

	return tmp.key, tmp.val
}

func (a *AVLTree[C]) Contains(key C) bool {
	return a.contains(a.root, key)
}

func (a *AVLTree[C]) contains(n *Node[C], key C) bool {
	if n == nil {
		return false
	}

	if n.key == key {
		return true
	}

	if key < n.key {
		return a.contains(n.left, key)
	}

	return a.contains(n.right, key)
}

func (a *AVLTree[C]) Size() uint32 {
	return a.nodeCount
}

func (a *AVLTree[C]) insert(n *Node[C], key C, val interface{}) *Node[C] {
	if n == nil {
		return &Node[C]{
			key: key,
			val: val,
		}
	}

	if key < n.key {
		n.left = a.insert(n.left, key, val)
	} else {
		n.right = a.insert(n.right, key, val)
	}

	a.update(n)

	return a.balance(n)
}

// updates a nodes height and balance factor
func (a *AVLTree[C]) update(n *Node[C]) {
	leftNodeHeight := -1
	if n.left != nil {
		leftNodeHeight = int(n.left.height)
	}

	rightNodeHeight := -1
	if n.right != nil {
		rightNodeHeight = int(n.right.height)
	}

	n.height = 1 + uint32(max(leftNodeHeight, rightNodeHeight))

	n.balanceFactor = rightNodeHeight - leftNodeHeight
}

func (a *AVLTree[C]) balance(n *Node[C]) *Node[C] {
	if n.balanceFactor == -2 {
		if n.left.balanceFactor <= 0 {
			return a.leftLeftCase(n)
		}

		return a.leftRightCase(n)
	}

	if n.balanceFactor == 2 {
		if n.right.balanceFactor >= 0 {
			return a.rightRightCase(n)
		}

		return a.rightLeftCase(n)
	}

	return n
}

func (a *AVLTree[C]) leftLeftCase(n *Node[C]) *Node[C] {
	return a.rightRotation(n)
}

func (a *AVLTree[C]) leftRightCase(n *Node[C]) *Node[C] {
	n.left = a.leftRotation(n.left)
	return a.leftLeftCase(n)
}

func (a *AVLTree[C]) rightRightCase(n *Node[C]) *Node[C] {
	return a.leftRotation(n)
}

func (a *AVLTree[C]) rightLeftCase(n *Node[C]) *Node[C] {
	n.right = a.rightRotation(n.right)
	return a.rightRightCase(n)
}

func (a *AVLTree[C]) leftRotation(n *Node[C]) *Node[C] {
	newParent := n.right
	n.right = newParent.left
	newParent.left = n

	a.update(n)
	a.update(newParent)

	return newParent
}

func (a *AVLTree[C]) rightRotation(n *Node[C]) *Node[C] {
	newParent := n.left
	n.left = newParent.right
	newParent.right = n

	a.update(n)
	a.update(newParent)

	return newParent
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
