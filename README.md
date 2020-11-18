# skewer

[![Build Status](https://travis-ci.org/sysread/skewer.svg?branch=master)](https://travis-ci.org/sysread/skewer)
[![codecov](https://codecov.io/gh/sysread/skewer/branch/master/graph/badge.svg)](https://codecov.io/gh/sysread/skewer)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/sysread/skewer)

Package skewer - a mergable priority queue

Skew heaps implement a priority queue (min heap) using a binary heap which
is continually rebalanced with each Put and Take operation.  Skew heaps have
an ammortized performance slighter better than O(log n).

The key feature of a skew heap is that it may be quickly and trivially
merged with another skew heap.  All heap operations are defined in terms of
the merge operation.

Mutable operations on the skew heap are atomic.

For more details, see [https://en.wikipedia.org/wiki/Skew_heap](https://en.wikipedia.org/wiki/Skew_heap)

## Examples

```golang
package main

import (
	"fmt"
	"github.com/sysread/skewer"
)

// Define a type that implements SkewItem. A SkewItem need only provide a
// single method, 'Priority', which returns the relative priority for an item
// in the queue. This value becomes the sorting mechanism for items in the
// heap. A lower value indicates a higher priority.
type Item int

func (item Item) Priority() int {
	// Negate the item's value so that a higher number will be given a higher
	// priority.
	return 0 - int(item)
}

func main() {
	heap := skewer.New()

	fmt.Println(heap.Top())

	for i := 0; i < 5; i++ {
		heap.Put(Item(i))
	}

	for i := 0; i < 5; i++ {
		fmt.Println(heap.Take())
	}

}

```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
