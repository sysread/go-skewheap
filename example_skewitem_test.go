package skewheap_test

import (
	"fmt"
	"github.com/sysread/skewheap"
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

func Example() {
	heap := skewheap.New()

	fmt.Println(heap.Top())

	for i := 0; i < 5; i++ {
		heap.Put(Item(i))
	}

	for i := 0; i < 5; i++ {
		fmt.Println(heap.Take())
	}

	// Output:
	// <nil> empty
	// 4 <nil>
	// 3 <nil>
	// 2 <nil>
	// 1 <nil>
	// 0 <nil>
}
