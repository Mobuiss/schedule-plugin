package filter

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Name is the name of the plugin used in the plugin registry and configurations.

const Name = "sample"

// Sort is a plugin that implements QoS class based sorting.

type sample struct{}

// 检验是否实现了framework.FilterPlugind的接口
var _ framework.FilterPlugin = &sample{}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &sample{}, nil
}

// Name returns name of the plugin.
func (pl *sample) Name() string {
	return Name
}

func (pl *sample) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	log.Printf("filter pod: %v, node: %v", pod.Name, nodeInfo)
	log.Println(state)
	// 排除没有user=dengdl24标签的节点
	if nodeInfo.Node().Labels["user"] != "dengdl24" {
		return framework.NewStatus(framework.Unschedulable, "Node: "+nodeInfo.Node().Name)
	}
	return framework.NewStatus(framework.Success, "Node: "+nodeInfo.Node().Name)
}
