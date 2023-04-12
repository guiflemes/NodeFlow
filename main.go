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
	json.Unmarshal([]byte(flowChartJson), dto)
	domain, _ := toDomain(dto)
	fmt.Println("resultDomain", domain)
}
