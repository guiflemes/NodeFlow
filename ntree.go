package main

import (
	"fmt"
)

const (
	TraverseInOrder TraverseType = iota
	TraversePreOrder
	TraversePostOrder
	TraverseLevelOrder
)

const (
	TraverseLeaves TraverseFlags = 1 << iota
	TraverseNonLeaves
	TraverseMask = 0x3
	TraverseAll  = TraverseLeaves | TraverseNonLeaves
)

type TraverseFunc[T comparable] func(*Node[T]) bool
type TraverseType int
type TraverseFlags int

type Node[T comparable] struct {
	data     T
	next     *Node[T]
	previous *Node[T]
	parent   *Node[T]
	children *Node[T]
}

func NewNode[T comparable](data T) *Node[T] {
	return &Node[T]{
		data: data,
	}
}

func (n *Node[T]) IsRoot() bool {
	return n.parent == nil && n.previous == nil && n.next == nil
}

func (n *Node[T]) Depth() int {
	depth := 0
	for n != nil {
		depth++
		n = n.parent
	}
	return depth
}

func (n *Node[T]) AddChild(child *Node[T]) *Node[T] {
	var node *Node[T]

	if !child.IsRoot() {
		return node
	}

	child.parent = n

	if n.children != nil {
		sibling := n.children
		for sibling.next != nil {
			sibling = sibling.next
		}
		child.previous = sibling
		sibling.next = child
		return node
	}

	child.parent.children = child
	return child
}

func (n *Node[T]) GetRoot() (*Node[T], int) {
	depth := 1
	current := n

	for current.parent != nil {
		depth++
		current = current.parent
	}

	return current, depth
}

func (n *Node[T]) traversePreOrder(flags TraverseFlags, traverseFunc TraverseFunc[T]) bool {

	if n.children != nil {
		func() bool {
			if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
				return true
			}

			child := n.children
			for child != nil {
				current := child
				child = current.next
				if current.traversePreOrder(flags, traverseFunc) {
					return true
				}
			}

			return false
		}()

	}

	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	return false
}

func (n *Node[T]) traverseInOrder(flags TraverseFlags, traverseFunc TraverseFunc[T]) bool {
	if n.children != nil {

		func() bool {
			child := n.children
			current := child
			child = current.next

			if current.traverseInOrder(flags, traverseFunc) {
				return true
			}

			if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
				return true
			}

			for child != nil {
				current = child
				child = current.next
				if current.traverseInOrder(flags, traverseFunc) {
					return true
				}
			}

			return false
		}()

	}

	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	return false
}

func (n *Node[T]) traversePostOrder(flags TraverseFlags, traverseFunc TraverseFunc[T]) bool {
	if n.children != nil {
		func() bool {
			child := n.children
			for child != nil {
				current := child
				child = current.next
				if current.traversePostOrder(flags, traverseFunc) {
					return true
				}
			}

			if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
				return true
			}

			return false
		}()
	}

	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	return false
}

func (n *Node[T]) depthTraversePreOrder(flags TraverseFlags, depth int, traverseFunc TraverseFunc[T]) bool {

	if n.children != nil {
		func() bool {

			if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
				return true
			}

			depth--
			if depth == 0 {
				return false
			}

			child := n.children

			for child != nil {
				current := child
				child = current.next

				if current.depthTraversePreOrder(flags, depth, traverseFunc) {
					return true
				}
			}

			return false
		}()
	}

	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	return false
}

func (n *Node[T]) depthTraverseInOrder(flags TraverseFlags, depth int, traverseFunc TraverseFunc[T]) bool {

	if n.children != nil {

		func() bool {
			depth--
			if depth > 0 {

				child := n.children
				current := child
				child = current.next

				if current.depthTraverseInOrder(flags, depth, traverseFunc) {
					return true
				}

				if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
					return true
				}

				for child != nil {
					current = child
					child = current.next
					if current.depthTraverseInOrder(flags, depth, traverseFunc) {
						return true
					}
				}

			}

			if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
				return true
			}

			return false
		}()

		return false
	}

	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	return false
}

func (n *Node[T]) depthTraversePostOrder(flags TraverseFlags, depth int, traverseFunc TraverseFunc[T]) bool {

	if n.children != nil {
		func() bool {
			depth--

			if depth > 0 {
				child := n.children

				for child != nil {
					current := child
					child = current.next
					if current.depthTraversePostOrder(flags, depth, traverseFunc) {
						return true
					}
				}
			}

			if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
				return true
			}
			return false
		}()
	}

	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	return false
}

func (n *Node[T]) Traverse(order TraverseType, flags TraverseFlags, depth int, trTraverseFunc TraverseFunc[T]) {

	if n == nil || trTraverseFunc == nil || order > TraverseLevelOrder || flags > TraverseMask || (depth < -1 || depth == 0) {
		return
	}

	switch order {

	default:
		fallthrough

	case TraversePreOrder:
		func() {
			if depth < 0 {
				n.traversePreOrder(flags, trTraverseFunc)
				return
			}
			n.depthTraversePreOrder(flags, depth, trTraverseFunc)
		}()

	case TraverseInOrder:
		func() {
			if depth < 0 {
				n.traverseInOrder(flags, trTraverseFunc)
				return
			}
			n.depthTraverseInOrder(flags, depth, trTraverseFunc)
		}()

	case TraversePostOrder:
		func() {
			if depth < 0 {
				n.traversePostOrder(flags, trTraverseFunc)
				return
			}
			n.depthTraversePostOrder(flags, depth, trTraverseFunc)
		}()

	}
}

func (n *Node[T]) String() string {
	if n == nil {
		return "()"
	}

	levels := []string{""}

	n.Traverse(TraversePreOrder, TraverseAll, -1, func(node *Node[T]) bool {
		currentLevel := 0
		nodeP := node.parent
		for nodeP != nil {
			currentLevel++
			if len(levels) <= currentLevel {
				levels = append(levels, "")
			}
			nodeP = nodeP.parent
		}
		levels[currentLevel] += fmt.Sprintf("(%v)", node.data) + "\t"
		return false

	})

	s := ""

	for _, v := range levels {
		s += v + "\n\n"
	}

	return s
}
