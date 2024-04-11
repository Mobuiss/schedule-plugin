package score

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	cpu_query = "100 - (avg by (instance) (irate(node_cpu_seconds_total{instance=\"%v\",mode=\"idle\"}[%v])) * 100)"
)

type PrometheusHandle struct {
	timeRange time.Duration
	address   string
	api       v1.API
}

func NewPrometheusHandle(timeRange time.Duration, address string) *PrometheusHandle {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		log.Printf("failed create prometheus client: %v\n", err)
	}
	return &PrometheusHandle{
		timeRange: timeRange,
		address:   address,
		api:       v1.NewAPI(client),
	}
}

func (p *PrometheusHandle) Get_CPU_Usage(nodeName string) (model.Value, error) {
	ctx := context.Background()
	query := fmt.Sprintf(cpu_query, nodeName, p.timeRange)
	result, warning, err := p.api.Query(ctx, query, time.Now())
	if len(warning) > 0 {
		log.Printf("Get CPU Usage Warning:%v\n", warning)
	}
	if err != nil {
		log.Printf("Get CPU Usage Err:%v\n", err)
	}

	return result, err
}
