// Package avl implements an easy-to-use AVL tree in Golang.
package avl

// Range is the basic node for the AVL tree leaves.
// The concept Range can be a real data range in most cases,  such as
// [LOWER_BOUND, UPPER_BOUND] .
// A Range can also be a single value, such as [VALUE_A, VALUE_A] .
type Range interface {
	// Compare returns an integer comparing two elements.
	// The result will be -1 if current element is less than the right,
	// and +1 if current element is greater that the right.
	// Otherwise, 0 will be return.
	Compare(right Range) int

	// Contains returns true if the right is a subset of current element.
	Contains(right Range) bool

	// Union returns the union of current element and the right.
	Union(right Range) Range
}

// Tree is a high-performance AVL tree.
type Tree struct {
	root *avlNode
}

// Insert a new Range into the AVL tree.
func (t *Tree) Insert(val Range) {
	t.root = t.root.insert(val)
}

// Search returns true if the AVL tree contains the <val>.
func (t *Tree) Search(val Range) bool {
	return t.root.search(val)
}

type avlNode struct {
	val Range

	parent *avlNode
	left   *avlNode
	right  *avlNode

	h int // the height
}

func (n *avlNode) height() int {
	if n == nil {
		return -1
	}
	return n.h
}

// updateHeight updates the height for current node and return the new height.
func (n *avlNode) updateHeight() int {
	n.h = n.left.height() + 1
	if rh := n.right.height() + 1; rh > n.h {
		n.h = rh
	}
	return n.h
}

// rotateLeft and other rotations are implemented according to the algorithm
// described on https://en.wikipedia.org/wiki/AVL_tree
func (n *avlNode) rotateLeft() (z *avlNode) {
	z, n.right = n.right, n.right.left
	z.left = n
	n.parent = z
	if n.right != nil {
		n.right.parent = n
	}

	n.updateHeight()
	z.updateHeight()
	return
}

func (n *avlNode) rotateRight() (z *avlNode) {
	z, n.left = n.left, n.left.right
	z.right = n
	n.parent = z
	if n.left != nil {
		n.left.parent = n
	}

	n.updateHeight()
	z.updateHeight()
	return
}

func (n *avlNode) rotateLeftRight() *avlNode {
	n.left = n.left.rotateLeft()
	return n.rotateRight()
}

func (n *avlNode) rotateRightLeft() *avlNode {
	n.right = n.right.rotateRight()
	return n.rotateLeft()
}

func (n *avlNode) rebalance() (p *avlNode) {
	for p = n.parent; p != nil; n, p = p, p.parent {
		grandParant, oldHeight := p.parent, p.height()
		leftChild := grandParant != nil && grandParant.left == p
		switch factor := p.left.height() - p.right.height(); {
		case factor > 1: // left heavy
			if n.left.height() < n.right.height() {
				p = p.rotateLeftRight()
			} else {
				p = p.rotateRight()
			}
		case factor < -1: // left heavy
			if n.left.height() > n.right.height() {
				p = p.rotateRightLeft()
			} else {
				p = p.rotateLeft()
			}
		}

		p.parent = grandParant
		if grandParant == nil {
			p.updateHeight()
			return
		}
		if leftChild {
			grandParant.left = p
		} else {
			grandParant.right = p
		}

		if p.updateHeight() == oldHeight {
			// the height for current node didn't change.
			break
		}
	}

	for p.parent != nil {
		p = p.parent
	}
	return
}

// insert a new <val> and return the new root of the AVL tree.
func (n *avlNode) insert(val Range) *avlNode {
	z := &avlNode{val: val}
	if n == nil {
		return z
	}

	for x := n; ; {
		switch factor := x.val.Compare(val); {
		case factor < 0: // x < z
			if x.right == nil {
				x.right, z.parent = z, x
				return z.rebalance()
			}
			x = x.right

		case factor > 0: // x > z
			if x.left == nil {
				x.left, z.parent = z, x
				return z.rebalance()
			}
			x = x.left

		default: // x == z
			if !x.val.Contains(val) {
				x.val = x.val.Union(val)
			}
			return n
		}
	}
}

// search returns true if the AVL tree contains the <val>.
func (n *avlNode) search(val Range) bool {
	for n != nil {
		switch factor := n.val.Compare(val); {
		case factor < 0: // n < z
			n = n.right
		case factor > 0: // n > z
			n = n.left
		default: // n == z
			return n.val.Contains(val)
		}
	}
	return false
}

// DebugPreorder will traverse the tree in preorder. For debug-use only.
func DebugPreorder(t *Tree) (ret []interface{}) {
	if t.root == nil {
		return nil
	}

	queue := []*avlNode{t.root}

	for idx := 0; idx < len(queue); idx++ {
		node := queue[idx]
		if node == nil {
			continue
		}
		// enqueue
		queue = append(queue, node.left, node.right)
		// fill the value
		ret = append(ret, node.val)
	}
	return
}
