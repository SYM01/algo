package avl

import "bytes"

type intRange int

func (i intRange) Compare(right Range) int {
	// The type of right must be intRange
	r := right.(intRange)
	switch {
	case i < r:
		return -1
	case i > r:
		return 1
	default:
		return 0
	}
}
func (i intRange) Contains(right Range) bool { return i == right }
func (i intRange) Union(right Range) Range   { return i }

// IntTree is a high-performance AVL tree for int.
type IntTree Tree

// Insert a new Range into the AVL tree.
func (t *IntTree) Insert(val int) {
	t.root = t.root.insert(intRange(val))
}

// Search returns true if the AVL tree contains the <val>.
func (t *IntTree) Search(val int) bool {
	return t.root.search(intRange(val))
}

type byteRange []byte

func (i byteRange) Compare(right Range) int   { return bytes.Compare(i, right.(byteRange)) }
func (i byteRange) Contains(right Range) bool { return bytes.Equal(i, right.(byteRange)) }
func (i byteRange) Union(right Range) Range   { return i }

// BytesTree is a high-performance AVL tree for []byte.
type BytesTree Tree

// Insert a new Range into the AVL tree.
func (t *BytesTree) Insert(val []byte) {
	t.root = t.root.insert(byteRange(val))
}

// Search returns true if the AVL tree contains the <val>.
func (t *BytesTree) Search(val []byte) bool {
	return t.root.search(byteRange(val))
}

// StringTree is a high-performance AVL tree for String.
type StringTree Tree

// Insert a new Range into the AVL tree.
func (t *StringTree) Insert(val string) {
	t.root = t.root.insert(byteRange(val))
}

// Search returns true if the AVL tree contains the <val>.
func (t *StringTree) Search(val string) bool {
	return t.root.search(byteRange(val))
}
