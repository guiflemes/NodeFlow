package main

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	Label string `json:"label"`
}

func (d *Data) String() string {
	return d.Label
}

func main() {
	dto := &FlowChartDto[Data]{}
	err := json.Unmarshal([]byte(flowChartJson), dto)
	fmt.Println(err)
	fmt.Println(dto.Nodes)
	domain, err := toDomain(dto)

	fmt.Println("errDomain", err)
	fmt.Println("resultDomain", domain)
}
