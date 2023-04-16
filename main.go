package main

import (
	"context"
	"encoding/json"
	"flowChart/adapters"
	"flowChart/domain"
	"flowChart/handlers"
	"flowChart/ports"
	"fmt"
)

func main() {
	config := &adapters.DatabaseConfig{}
	config.Parse()
	repo := adapters.NewFlowChartDataRepo(config)

	dto := &ports.FlowChartDto[ports.DataDto]{}
	json.Unmarshal([]byte(ports.FlowChartJson), dto)
	domain, _ := handlers.ToDomain(dto, func(request ports.DataDto) domain.Data {
		return domain.Data{Label: request.Label}
	})

	err := repo.StoreFlowChart(context.Background(), domain)
	fmt.Println(err)

}
