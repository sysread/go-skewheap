package skewheap_test

import "fmt"
import "math/rand"
import "sort"
import "testing"
import "github.com/sysread/skewheap"

type IntItem int

func (item IntItem) Priority() int {
	return int(item)
}

func is(t *testing.T, err error, got int, expected int, msg string) bool {
	if err != nil {
		t.Log("FAIL", msg)
		return false
	} else if got != expected {
		t.Log("FAIL", msg)
		t.Log("expected:", expected)
		t.Log("  actual:", got)
		t.Fail()
		return false
	}

	return true
}

func TestPut(t *testing.T) {
	heap := skewheap.New()
	is(t, nil, heap.Size(), 0, "initial heap Size()")

	heap.Put(IntItem(42))
	is(t, nil, heap.Size(), 1, "put 1")

	heap.Put(IntItem(10))
	is(t, nil, heap.Size(), 2, "put 2")
}

func TestTake(t *testing.T) {
	heap := skewheap.New()

	ints := rand.Perm(50)

	for _, i := range ints {
		heap.Put(IntItem(i))
	}

	sort.Sort(sort.IntSlice(ints))

	for _, i := range ints {
		top, err1 := heap.Top()
		is(t, err1, int(top.(IntItem)), i, fmt.Sprintf("Top() == %d", i))

		val, err2 := heap.Take()
		is(t, err2, int(val.(IntItem)), i, fmt.Sprintf("Take() == %d", i))
	}

	top, err1 := heap.Top()

	if top != nil {
		t.Log("Top() did not return nil when called from empty heap")
	}

	if fmt.Sprintf("%v", err1) != "empty" {
		t.Log("Top() did not return expected error when called from empty heap")
	}

	val, err2 := heap.Take()

	if val != nil {
		t.Log("Take() did not return nil when called from empty heap")
	}

	if fmt.Sprintf("%v", err2) != "empty" {
		t.Log("Take() did not return expected error when called from empty heap")
	}
}

func TestMerge(t *testing.T) {
	aInts := []int{0, 1, 2, 3, 4}
	bInts := []int{5, 6, 7, 8, 9}
	cInts := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	a, b := skewheap.New(), skewheap.New()

	for _, i := range aInts {
		a.Put(IntItem(i))
	}

	for _, i := range bInts {
		b.Put(IntItem(i))
	}

	c := a.Merge(*b)

	is(t, nil, b.Size(), 5, "b.Size() remains intact")
	for _, i := range bInts {
		top, err1 := b.Top()
		is(t, err1, int(top.(IntItem)), i, fmt.Sprintf("b.Top() == %d", i))

		val, err2 := b.Take()
		is(t, err2, int(val.(IntItem)), i, fmt.Sprintf("b.Take() == %d", i))
	}

	is(t, nil, c.Size(), 10, "c.Size() is a.Size() + b.Size()")
	for _, i := range cInts {
		top, err1 := c.Top()
		is(t, err1, int(top.(IntItem)), i, fmt.Sprintf("c.Top() == %d", i))

		val, err2 := c.Take()
		is(t, err2, int(val.(IntItem)), i, fmt.Sprintf("c.Take() == %d", i))
	}

	is(t, nil, a.Size(), 5, "a.Size() remains intact")
	for _, i := range aInts {
		top, err1 := a.Top()
		is(t, err1, int(top.(IntItem)), i, fmt.Sprintf("a.Top() == %d", i))

		val, err2 := a.Take()
		is(t, err2, int(val.(IntItem)), i, fmt.Sprintf("a.Take() == %d", i))
	}
}
