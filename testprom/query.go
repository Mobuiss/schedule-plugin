package main

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func init() {
	var timerange time.Duration
	timerange = time.Second * 70
	fmt.Println(timerange)
	fmt.Println(fmt.Sprintf("100 - (avg by (instance) (irate(node_cpu_seconds_total{instance=\"%v\",mode=\"idle\"}[%v])) * 100)", "exporter", timerange))
}
func main() {
	fmt.Println("hello world")
	client, err := api.NewClient(api.Config{
		Address: "http://192.168.27.128:9090",
	})
	if err != nil {
		fmt.Println("???")
	}
	api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	const (
		ask1 = "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"
		ask2 = "node_load1{}"
		ask3 = "sum by (instance) (irate(node_network_transmit_bytes_total{device!~\"bond.*?|lo\"}[5m])/128)"
	)
	result, warning, err := api.Query(ctx, ask1, time.Now())
	if err != nil {
		fmt.Println(err)
	}
	if len(warning) > 0 {
		fmt.Println(warning)
	}
	vector := result.(model.Vector)
	fmt.Println(vector)
	fmt.Printf("%v\n", vector[0].Timestamp)
	for _, x := range vector {
		fmt.Printf("x=%v\n", x)
	}
	fmt.Printf("result:%v\n", result)
}
