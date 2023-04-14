package domain

type FlowChart[T comparable] struct {
	Id    string
	Title string
	Node  *Node[T]
}
