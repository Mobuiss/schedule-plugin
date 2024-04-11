package score

import (
	"context"
	"fmt"
	"time"

	"github.com/emicklei/go-restful/log"
	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/scheduler-plugins/apis/config"
)

var _ framework.ScorePlugin = &myScore{}

type myScore struct {
	prometheus *PrometheusHandle
	handle     framework.Handle
}

const Name = "score"

func (ms *myScore) Name() string {
	return Name
}

func New(obj runtime.Object, fh framework.Handle) (framework.Plugin, error) {

	args, ok := obj.(*config.MyScoreArgs)
	if !ok {
		log.Printf("Expect MyScoreArgs,got %T:%v\n", args, args)
		args = &config.MyScoreArgs{
			Address:      "http://192.168.27.128:9090",
			TimeInSecond: 60,
		}
	}
	return &myScore{
		prometheus: NewPrometheusHandle(
			time.Duration(args.TimeInSecond)*time.Second,
			args.Address),
		handle: fh,
	}, nil
}
func (ms *myScore) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	log.Printf("nodeName=%v\n", nodeName)
	result, err := ms.prometheus.Get_CPU_Usage(nodeName)
	if err != nil {
		log.Printf("Query CPU Usage Error:%v", err)
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Query CPU Usage Error:%v", err))
	}
	res := result.(model.Vector)
	log.Printf("node:%v,score:%v\n", nodeName, res[0].Value)
	return int64(res[0].Value * 1000), framework.NewStatus(framework.Success, fmt.Sprintf("success score node:%v\n", nodeName))
}

func (ms *myScore) ScoreExtensions() framework.ScoreExtensions {
	return ms
}

func (ms *myScore) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var maxn int64
	for _, x := range scores {
		if x.Score > maxn {
			maxn = x.Score
		}
	}
	for i, x := range scores {
		scores[i].Score = x.Score * framework.MaxNodeScore / maxn
	}
	log.Printf("node final scores:\n%v\n", scores)
	return nil
}
