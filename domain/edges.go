package domain

type Edge struct {
	Id     string
	Source string
	Target string
}

func NewEdge(id string, source string, target string) *Edge {
	return &Edge{
		Id:     id,
		Source: source,
		Target: target,
	}
}
