package domain

type FlowChart[T comparable] struct {
	id   string
	Node *Node[T]
}

// func NewFlowChart[T comparable](Node *Node[T]) *FlowChart[T] {
// 	return &FlowChart[T]{
// 		id: "",
// 		n:  Node,
// 	}
// }

// func (f *FlowChart[T]) String() string {
// 	return fmt.Sprintf("id=%s, nodesCount=(%d)", f.id, 1)
// }
