package main

import (
	"context"
	"encoding/json"
	"flowChart/adapters"
	"flowChart/handlers"
	"flowChart/ports"
	"fmt"
)

func main() {
	config := &adapters.DatabaseConfig{}
	config.Parse()
	repo := adapters.NewFlowChartDataRepo(config)

	dto := &ports.FlowChartDto[ports.Data]{}
	json.Unmarshal([]byte(ports.FlowChartJson), dto)
	domain, _ := handlers.ToDomain(dto)

	err := repo.StoreFlowChart(context.Background(), domain)
	fmt.Println(err)

}
