//go:build go1.18
// +build go1.18

package avl_test

import (
	"fmt"

	"github.com/sym01/algo/avl"
)

func ExampleNewOrderedTree_int() {
	ints := []int{
		9, 85, 45, 76, 53, 52, 66, 12, 13, 65,
		98, 33, 28, 42, 38, 84, 24, 37, 36, 35,
		27, 41, 77, 48, 46, 81, 78, 60, 25, 5,
		3, 93, 11, 49, 74, 96, 94, 39, 51, 95,
		58, 90, 91, 20, 7, 44, 97, 29, 2, 8,
		19, 14, 86, 61, 54, 79, 64, 16, 67, 6,
		26, 73, 50, 72, 10, 83, 70, 80, 89, 15,
		17, 4, 100, 71, 55, 75, 69, 87, 34, 92,
		43, 18, 88, 82, 30, 57, 23, 99, 59, 21,
		32, 22, 68, 56, 47, 40, 63, 62, 1, 31,
	}

	tree := avl.NewOrderedTree[int]()
	for _, i := range ints {
		tree.Insert(i)
	}

	fmt.Println(tree.Search(42))
	fmt.Println(tree.Search(96))
	fmt.Println(tree.Search(1024))

	// Output:
	// true
	// true
	// false
}

func ExampleNewOrderedTree_string() {
	strs := []string{
		"eKSI9wZlde",
		"j1ORx3W0ph",
		"1Lw7uwm4VS",
		"WcOTlexQqG",
		"hqDGapsHq1",
		"mDy5mwwNzx",
		"VhY3HfqmG5",
		"Y5doxlrZe6",
		"VbR2pKkv9E",
		"OpUWxSI2eP",
		"Ja6CqjnFq5",
		"h6dlJ5m",
		"PBBKZnRBYu",
		"AgX3njkaZT",
		"ytXkDnD0vr",
		"8QJFS2fLGd",
		"PAarEfdt1w", // same string
		"PAarEfdt1w", // same string
		"PAarEfdt1w", // same string
		"8QJFS2fLGd",
	}

	tree := avl.NewOrderedTree[string]()
	for _, s := range strs {
		tree.Insert(s)
	}

	for _, s := range strs {
		if !tree.Search(s) {
			fmt.Printf("%s not in tree\n", s)
		}
	}

	s := "abccc"
	if !tree.Search(s) {
		fmt.Printf("%s not in tree\n", s)
	}

	// Output:
	// abccc not in tree
}
