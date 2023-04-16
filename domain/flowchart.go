package domain

type FlowChart[T comparable] struct {
	Id    string
	Title string
	Key   string
	Node  *Node[T]
}
