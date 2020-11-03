package avl_test

import (
	"fmt"
	"testing"

	"github.com/sym01/algo/avl"
)

func TestAvlNode_rebalance(t *testing.T) {
	tree := new(avl.IntTree)
	if avl.DebugPreorder((*avl.Tree)(tree)) != nil {
		t.Error("unexpected behavior of DebugPreorder")
	}

	for i := 1; i < 11; i++ {
		tree.Insert(i)
	}

	expected := []int{4, 2, 8, 1, 3, 6, 9, 5, 7, 10}
	ret := avl.DebugPreorder((*avl.Tree)(tree))

	for i := 0; i < len(expected); i++ {
		if fmt.Sprint(ret[i]) != fmt.Sprint(expected[i]) {
			t.Fatalf("unexpected result, expect %v, got %v", expected, ret)
		}
	}
}

func TestBytesTree(t *testing.T) {
	tree := new(avl.BytesTree)

	for i := byte(1); i < 20; i++ {
		tree.Insert([]byte{i, i, i})
	}

	testcases := []struct {
		bytes    []byte
		expected bool
	}{
		{[]byte("\x10\x10\x10"), true},
		{[]byte("\x10\x12\x10"), false},
	}

	for _, item := range testcases {
		if ret := tree.Search(item.bytes); ret != item.expected {
			t.Fatalf("unexpected result, expect %v for %v, got %v",
				item.expected, item.bytes, ret)
		}
	}

}
