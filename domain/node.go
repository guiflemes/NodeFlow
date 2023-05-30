package domain

import (
	"database/sql/driver"
	"encoding/json"
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

type TraverseFunc[T any] func(*Node[T]) bool
type TraverseType int
type TraverseFlags int

type Position struct {
	X float64
	Y float64
}

func (p Position) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type Node[T any] struct {
	NodeID           string
	Data             T
	Position         Position
	Width            int16
	Height           int16
	Selected         bool
	PositionAbsolute Position
	Dragging         bool
	Type             string
	next             *Node[T]
	previous         *Node[T]
	parent           *Node[T]
	children         *Node[T]
}

func NewNode[T any](nodeID string, data T, position Position, width int16, height int16, selected bool, positionAbsolute Position, dragging bool, typeNode string) *Node[T] {
	return &Node[T]{
		NodeID:           nodeID,
		Data:             data,
		Position:         position,
		Width:            width,
		Height:           height,
		Selected:         selected,
		PositionAbsolute: positionAbsolute,
		Dragging:         dragging,
		Type:             typeNode,
	}
}

func (n *Node[T]) ParentId() string {
	if n.IsRoot() {
		return "0"
	}
	return n.parent.NodeID
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
	var Node *Node[T]

	if !child.IsRoot() {
		return Node
	}

	child.parent = n

	if n.children != nil {
		sibling := n.children
		for sibling.next != nil {
			sibling = sibling.next
		}
		child.previous = sibling
		sibling.next = child
	} else {
		child.parent.children = child
	}

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
	// First, apply the function to the current node if the flag is set
	if flags&TraverseNonLeaves != 0 && traverseFunc(n) {
		return true
	}

	// Traverse the children recursively
	child := n.children
	for child != nil {
		if child.traversePreOrder(flags, traverseFunc) {
			return true
		}
		child = child.next
	}

	// If the flag is set, apply the function to the current node again
	if flags&TraverseNonLeaves == 0 && traverseFunc(n) {
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

	n.Traverse(TraversePreOrder, TraverseAll, -1, func(Node *Node[T]) bool {
		currentLevel := 0
		nodeP := Node.parent
		for nodeP != nil {
			currentLevel++
			if len(levels) <= currentLevel {
				levels = append(levels, "")
			}
			nodeP = nodeP.parent
		}
		levels[currentLevel] += fmt.Sprintf("NodeID=%s data=(%v)", Node.NodeID, Node.Data) + "\t"
		return false

	})

	s := ""

	for _, v := range levels {
		s += v + "\n\n"
	}

	return s
}
