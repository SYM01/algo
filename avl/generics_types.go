//go:build go1.18
// +build go1.18

package avl

import (
	"golang.org/x/exp/constraints"
)

// ITree is an AVL tree implement with type parameters support.
type ITree[T any] interface {
	Insert(T)
	Search(T) bool
}

// NewOrderedTree creates a new high-performance AVL tree instance for
// ordered types, such as integer, float, and string.
func NewOrderedTree[T constraints.Ordered]() ITree[T] {
	return new(orderedTree[T])
}

type orderedRange[T constraints.Ordered] struct {
	v T
}

func (i orderedRange[T]) Compare(right Range) int {
	// The type of right must be intRange
	r := right.(orderedRange[T])

	switch {
	case i.v < r.v:
		return -1
	case i.v > r.v:
		return 1
	default:
		return 0
	}
}
func (i orderedRange[T]) Contains(right Range) bool { return i.v == right.(orderedRange[T]).v }
func (i orderedRange[T]) Union(right Range) Range   { return i }

type orderedTree[T constraints.Ordered] struct {
	Tree
}

func (i *orderedTree[T]) Insert(v T) {
	i.Tree.Insert(orderedRange[T]{v})
}
func (i *orderedTree[T]) Search(v T) bool {
	return i.Tree.Search(orderedRange[T]{v})
}
