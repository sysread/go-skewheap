# skewheap
--
    import "github.com/sysread/skewheap"

Skew heaps implement a priority queue (min heap) using a binary heap which is
continually rebalanced with each Put and Take operation. Skew heaps have an
ammortized performance slighter better than O(log n).

The key feature of a skew heap is that it may be quickly and trivially merged
with another skew heap. All heap operations are defined in terms of the merge
operation.

For more details, see https://en.wikipedia.org/wiki/Skew_heap

## Usage

#### type SkewHeap

```go
type SkewHeap struct {
	// The number of items in the queue.
	Size int
}
```

SkewHeap is the base interface type. It's only exposed member is Size.

#### func  New

```go
func New() *SkewHeap
```

#### func (SkewHeap) Explain

```go
func (heap SkewHeap) Explain()
```
Debugging routine that emits a description of the skew heap and its internal
structure to stdout.

#### func (SkewHeap) Merge

```go
func (heap SkewHeap) Merge(other SkewHeap) *SkewHeap
```
Non-destructively combines two heaps into a new heap. Note that Merge
recursively copies the structure of each input heap.

#### func (*SkewHeap) Put

```go
func (heap *SkewHeap) Put(value SkewItem)
```
Inserts a value into the heap.

#### func (*SkewHeap) Take

```go
func (heap *SkewHeap) Take() (SkewItem, error)
```
Removes and returns the value with the highest priority from the heap.

#### func (*SkewHeap) Top

```go
func (heap *SkewHeap) Top() (SkewItem, error)
```
Returns the value highest priority from the heap without removing it.

#### type SkewItem

```go
type SkewItem interface {
	Priority() int
}
```

The skew heap can queue any item that can provide a relative priority value by
implementing the Priority() method. A lower value indicates a higher priority in
the queue.
