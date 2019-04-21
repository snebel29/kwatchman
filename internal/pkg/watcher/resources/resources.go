package resources

import (
	"fmt"
	_ "github.com/snebel29/kooper/operator/common"
	_ "github.com/snebel29/kooper/operator/controller"
	_ "github.com/snebel29/kooper/operator/handler"
	_ "github.com/snebel29/kooper/operator/retrieve"
)

type ResourceWatcher interface {
	Run() error
	Shutdown()
}

type K8sResourceWatcher struct {
	kind string
}

func (r *K8sResourceWatcher) Run() error {
	fmt.Printf("Run K8sResourceWatcher with kind %v", r.kind)
	return nil
}

func (r *K8sResourceWatcher) Shutdown() {
	fmt.Printf("Shutdown K8sResourceWatcher with kind %v", r.kind)
}

func NewK8sDeploymentWatcher() ResourceWatcher {
	return &K8sResourceWatcher{kind: "Deployment"}
}
