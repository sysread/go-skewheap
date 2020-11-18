// Package skewer - a mergable priority queue
//
// Skew heaps implement a priority queue (min heap) using a binary heap which
// is continually rebalanced with each Put and Take operation.  Skew heaps have
// an ammortized performance slighter better than O(log n).
//
// The key feature of a skew heap is that it may be quickly and trivially
// merged with another skew heap.  All heap operations are defined in terms of
// the merge operation.
//
// Mutable operations on the skew heap are atomic.
//
// For more details, see https://en.wikipedia.org/wiki/Skew_heap
package skewer

import "errors"
import "fmt"
import "sort"
import "sync"

// SkewItem requires any queueable value to be wrapped in an interface which
// can provide a relative priority for the value as an integer. A lower value
// indicates a higher priority in the queue.
type SkewItem interface {
	Priority() int
}

type skewNode struct {
	left, right *skewNode
	value       SkewItem
}

func (node skewNode) priority() int {
	return node.value.Priority()
}

// SkewHeap is the base interface type
type SkewHeap struct {
	// The number of items in the queue
	size  int
	mutex *sync.Mutex
	root  *skewNode
}

// Size returns the number of items in the queue.
func (heap SkewHeap) Size() int { return heap.size }

// Sort interface
type byPriority []*skewNode

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority() < a[j].priority() }

// New initializes and returns a new *SkewHeap.
func New() *SkewHeap {
	heap := &SkewHeap{
		size:  0,
		mutex: &sync.Mutex{},
		root:  nil,
	}

	return heap
}

// Voluntarily locks the data structure while modifying it.
func (heap *SkewHeap) lock()   { heap.mutex.Lock() }
func (heap *SkewHeap) unlock() { heap.mutex.Unlock() }

// Indents explain()
func indent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("    ")
	}
}

// Debugging routine that emits a description of the node and its internal
// structure to stdout.
func (node skewNode) explain(depth int) {
	indent(depth)
	fmt.Printf("Node<value:%v, priority:%d>\n", node.value, node.priority())

	if node.left != nil {
		indent(depth)
		fmt.Printf("-Left:\n")
		node.left.explain(depth + 1)
	}

	if node.right != nil {
		indent(depth)
		fmt.Printf("-Right:\n")
		node.right.explain(depth + 1)
	}
}

// Explain emits a description of the skew heap and its internal structure to
// stdout.
func (heap SkewHeap) Explain() {
	fmt.Printf("Heap<Size:%d>\n", heap.Size())
	fmt.Printf("-Root:\n")

	if heap.Size() > 0 {
		heap.root.explain(1)
	}

	fmt.Printf("\n")
}

// Merges two nodes destructively
func (heap *skewNode) merge(other *skewNode) *skewNode {
	if heap == nil {
		return other
	}

	if other == nil {
		return heap
	}

	// Cut the right subtree from each path and store the remaining left subtrees
	// in nodes.
	todo := [](*skewNode){heap, other}
	nodes := [](*skewNode){}

	for len(todo) > 0 {
		node := todo[0]
		todo = todo[1:]

		if node.right != nil {
			todo = append(todo, node.right)
			node.right = nil
		}

		nodes = append(nodes, node)
	}

	// Sort the cut paths
	sort.Sort(byPriority(nodes))

	// Recombine subtrees
	var node *skewNode

	for len(nodes) > 1 {
		node, nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
		prev := nodes[len(nodes)-1]

		// Set penultimate node's right child to its left (and only) subtree
		prev.right = prev.left

		// Set its left child to the ultimate node
		prev.left = node
	}

	return nodes[0]
}

// Recursively copies a node and its children
func (src *skewNode) copyNode() *skewNode {
	if src == nil {
		return nil
	}

	newNode := &skewNode{
		value: src.value,
		left:  src.left.copyNode(),
		right: src.right.copyNode(),
	}

	return newNode
}

// Merge non-destructively combines two heaps into a new heap. Note that Merge
// recursively copies the structure of each input heap.
func (heap SkewHeap) Merge(other SkewHeap) *SkewHeap {
	ready := make(chan bool, 2)

	var rootA, rootB *skewNode
	var sizeA, sizeB int

	// Because each heap may be used by other go routines, locking their mutexes
	// and copying their contents is done in another routine, and this thread
	// blocks on receiving a signal from the locking thread. This helps to avoid
	// unnecessary blocking by attempting to lock two mutexes serially.

	go func() {
		heap.lock()
		sizeA = heap.Size()
		rootA = heap.root.copyNode()
		heap.unlock()
		ready <- true
	}()

	go func() {
		other.lock()
		sizeB = other.Size()
		rootB = other.root.copyNode()
		other.unlock()
		ready <- true
	}()

	// Wait on copies to be made
	<-ready
	<-ready

	newHeap := New()
	newHeap.size += sizeA + sizeB
	newHeap.root = rootA.merge(rootB)

	return newHeap
}

// Put inserts a value into the heap.
func (heap *SkewHeap) Put(value SkewItem) {
	newNode := &skewNode{
		left:  nil,
		right: nil,
		value: value,
	}

	heap.lock()

	if heap.Size() == 0 {
		heap.root = newNode
	} else {
		heap.root = heap.root.merge(newNode)
	}

	heap.size++

	heap.unlock()
}

// Take removes and returns the value with the highest priority from the heap.
func (heap *SkewHeap) Take() (SkewItem, error) {
	heap.lock()

	if heap.Size() > 0 {
		value := heap.root.value
		heap.root = heap.root.left.merge(heap.root.right)
		heap.size--
		heap.unlock()
		return value, nil
	}

	heap.unlock()
	return nil, errors.New("empty")
}

// Top returns the value highest priority from the heap without removing it.
func (heap *SkewHeap) Top() (SkewItem, error) {
	if heap.Size() > 0 {
		return heap.root.value, nil
	}

	return nil, errors.New("empty")
}
