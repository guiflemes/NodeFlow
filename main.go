package main

import (
	"context"
	"encoding/json"
	"flowChart/repository"
	"fmt"
	"time"
)

type Data struct {
	Label string `json:"label"`
}

func (d *Data) String() string {
	return d.Label
}

func main() {
	config := &repository.Database{}
	config.Parse()
	repo := repository.NewPostgresRepo[Data](config)

	dto := &FlowChartDto[Data]{}
	json.Unmarshal([]byte(flowChartJson), dto)
	domain, _ := toDomain(dto)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := repo.Store(ctx, domain)
	fmt.Println(err)

}
