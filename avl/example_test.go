package avl_test

import (
	"fmt"

	"github.com/sym01/algo/avl"
)

func ExampleIntTree() {
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

	tree := new(avl.IntTree)
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

func ExampleStringTree() {
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

	tree := new(avl.StringTree)
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

// customized data struct for AVL tree
type intRange struct {
	min int
	max int
}

func (l *intRange) Compare(right avl.Range) int {
	r := right.(*intRange)
	if l.max < r.min {
		return -1
	}
	if l.min > r.max {
		return 1
	}
	return 0
}

func (l *intRange) Contains(right avl.Range) bool {
	r := right.(*intRange)
	if l.min > r.min {
		return false
	}
	if l.max < r.max {
		return false
	}

	return true
}

func (l *intRange) Union(right avl.Range) avl.Range {
	r := right.(*intRange)
	ret := &intRange{
		min: l.min,
		max: l.max,
	}
	if ret.min > r.min {
		ret.min = r.min
	}
	if ret.max < r.max {
		ret.max = r.max
	}
	return ret
}

func ExampleTree() {
	data := []*intRange{
		{10, 15},
		{20, 25},
		{30, 35},
		{40, 45},
		{21, 26},
		{50, 55},
		{60, 65},
	}
	tree := new(avl.Tree)
	for _, item := range data {
		tree.Insert(item)
	}

	fmt.Println(tree.Search(&intRange{20, 26}))
	fmt.Println(tree.Search(&intRange{20, 27}))

	// Output:
	// true
	// false
}
