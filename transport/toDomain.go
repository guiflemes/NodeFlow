package transport

import (
	"errors"
	"flowChart/domain"
	"sync"
)

func ToDomain[R comparable, D comparable](flowChart *FlowChartDto[R], dataParse func(request R) D) (*domain.FlowChart[D], error) {
	m := sync.RWMutex{}

	nodeMap := func() map[string]*domain.Node[D] {
		nodeM := make(map[string]*domain.Node[D], len(flowChart.Nodes))
		for _, n := range flowChart.Nodes {
			m.Lock()
			data := dataParse(n.Data)
			nodeM[n.Id] = domain.NewNode(n.Id, data, domain.Position{X: n.Position.X, Y: n.Position.Y},
				n.Width, n.Height, n.Selected, domain.Position{X: n.PositionAbsolute.X, Y: n.PositionAbsolute.Y}, n.Dragging)
			m.Unlock()
		}
		return nodeM
	}()

	var flow *domain.FlowChart[D]
	var root *domain.Node[D]

	for _, edge := range flowChart.Edges {

		if err := func() error {
			m.RLock()
			defer m.RUnlock()
			parent, ok := nodeMap[edge.Source]

			if !ok {
				return errors.New("Parent Not Found")
			}

			if edge.Source == "0" {
				root = parent
			}

			if child, ok := nodeMap[edge.Id]; ok {
				parent.AddChild(child)
				return nil
			}

			return errors.New("Child Not Found")

		}(); err != nil {
			return flow, err
		}

	}

	return &domain.FlowChart[D]{Title: flowChart.Title, Node: root, Key: flowChart.Key}, nil
}
