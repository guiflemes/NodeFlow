package domain

type FlowChart[T any] struct {
	Id    string
	Title string
	Key   string
	Node  *Node[T]
}
