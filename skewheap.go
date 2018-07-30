// Skew heaps implement a priority queue (min heap) using a binary heap which
// is continually rebalanced with each Put and Take operation. Skew heaps have
// an ammortized performance slighter better than O(log n).
//
// For more details, see https://en.wikipedia.org/wiki/Skew_heap
package skewheap

import "errors"
import "fmt"
import "sort"

// The skew heap can queue any item that can provide a relative priority value
// by implementing the Priority() method. A lower value indicates a higher
// priority in the queue.
type SkewItem interface {
	Priority() int
}

type node struct {
	left, right *node
	value       SkewItem
}

func (n node) priority() int { return n.value.Priority() }

// SkewHeap is the base interface type. It's only exposed member is Size.
type SkewHeap struct {
	// The number of items in the queue.
	Size int
	root *node
	sem  chan bool
}

// Sort interface
type byPriority []*node

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority() < a[j].priority() }

func New() *SkewHeap {
	heap := &SkewHeap{Size: 0, root: nil, sem: make(chan bool, 1)}
	heap.unlock()
	return heap
}

// Voluntarily locks the data structure while modifying it.
func (heap *SkewHeap) lock()   { <-heap.sem }
func (heap *SkewHeap) unlock() { heap.sem <- true }

// Indents explain()
func indent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("    ")
	}
}

// Debugging routine that emits a description of the node and its internal
// structure to stdout.
func (node node) explain(depth int) {
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

// Debugging routine that emits a description of the skew heap and its internal
// structure to stdout.
func (heap SkewHeap) Explain() {
	fmt.Printf("Heap<Size:%d>\n", heap.Size)
	fmt.Printf("-Root:\n")

	if heap.Size > 0 {
		heap.root.explain(1)
	}

	fmt.Printf("\n")
}

// Merges two nodes destructively
func (heap *node) merge(other *node) *node {
	if heap == nil {
		return other
	}

	if other == nil {
		return heap
	}

	// Cut the right subtree from each path and store the remaining left subtrees
	// in nodes.
	todo := []*node{heap, other}
	nodes := []*node{}

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
	var node *node

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
func (src *node) copyNode() *node {
	if src == nil {
		return nil
	}

	newNode := &node{
		value: src.value,
		left:  src.left.copyNode(),
		right: src.right.copyNode(),
	}

	return newNode
}

// Non-destructively combines two heaps into a new heap. Note that Merge
// recursively copies the structure of each input heap.
func (heap SkewHeap) Merge(other SkewHeap) *SkewHeap {
	ready := make(chan bool, 2)

	var rootA, rootB *node
	var sizeA, sizeB int

	go func() {
		heap.lock()
		sizeA = heap.Size
		rootA = heap.root.copyNode()
		heap.unlock()
		ready <- true
	}()

	go func() {
		other.lock()
		sizeB = other.Size
		rootB = other.root.copyNode()
		other.unlock()
		ready <- true
	}()

	<-ready
	<-ready

	newHeap := New()
	newHeap.Size += sizeA + sizeB
	newHeap.root = rootA.merge(rootB)

	return newHeap
}

// Inserts a value into the heap.
func (heap *SkewHeap) Put(value SkewItem) {
	newNode := &node{
		left:  nil,
		right: nil,
		value: value,
	}

	heap.lock()

	if heap.Size == 0 {
		heap.root = newNode
	} else {
		heap.root = heap.root.merge(newNode)
	}

	heap.Size++

	heap.unlock()
}

// Removes and returns the value with the highest priority from the heap.
func (heap *SkewHeap) Take() (SkewItem, error) {
	heap.lock()

	if heap.Size > 0 {
		value := heap.root.value
		heap.root = heap.root.left.merge(heap.root.right)
		heap.Size--
		heap.unlock()
		return value, nil
	} else {
		heap.unlock()
		return nil, errors.New("empty")
	}
}

// Returns the value highest priority from the heap without removing it.
func (heap *SkewHeap) Top() (SkewItem, error) {
	if heap.Size > 0 {
		return heap.root.value, nil
	} else {
		return nil, errors.New("empty")
	}
}
