package resources

import (
	"fmt"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
)

var thisFilename string

func init() {
	_, t, _, _ := runtime.Caller(0)
	thisFilename = t
}

type fakeDeployment struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

func TestMarshal(t *testing.T) {
	r, err := marshal(&fakeDeployment{Kind: "Deployment", ApiVersion: "extensions/v1beta1"})
	expected := "{\"kind\":\"Deployment\",\"apiVersion\":\"extensions/v1beta1\"}"
	if string(r) != expected {
		t.Errorf("%s Should match with %s", string(r), expected)
	}
	if err != nil {
		t.Error("err should be nil")
	}
	//It's hard to make Marshal to return an error, pretty much anything can be marshaled
	r, err = marshal(make(chan int))
	if err == nil {
		t.Error("Error was expected")
	}
}

func NewFakeDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		metav1.ObjectMeta{},
		appsv1.DeploymentSpec{},
		appsv1.DeploymentStatus{},
	}
}

func TestGetManifest(t *testing.T) {

	r, err := getManifest(NewFakeDeployment())
	if err != nil {
		t.Error("err should be nil")
	}
	expected := `{"kind":"Deployment","apiVersion":"extensions/v1beta1","metadata":{"creationTimestamp":null},"spec":{"selector":null,"template":{"metadata":{"creationTimestamp":null},"spec":{"containers":null}},"strategy":{}},"status":{}}`
	if string(r) != expected {
		t.Errorf("%s Should match with %s", string(r), expected)
	}

	r, err = getManifest(&appsv1.ReplicaSet{})
	if err == nil {
		t.Error("err should be error")
	}
}

func TestNewKooperHandlerFunctionWithAdd(t *testing.T) {
	h1 := handler.NewMockHandler()
	chainOfHandlers := handler.NewChainOfHandlers(h1)
	fn := newKooperHandlerFunction(chainOfHandlers, "Deployment")

	err := fn(nil, &common.K8sEvent{
		Kind:      "Add",
		HasSynced: true,
		Key:       "default/den-from-neverwhere",
		Object:    NewFakeDeployment(),
	})
	if err != nil {
		t.Error("err should be nil")
	}

	if h1.Called != true {
		t.Error("Handler should have been called")
	}
	if reflect.DeepEqual(h1.PassedK8sManifest, []byte{}) {
		t.Error("No K8sManifest have been set")
	}
	if !reflect.DeepEqual(h1.PassedPayload, []byte{}) {
		t.Error("Payload should be empty")
	}
}

func TestNewKooperHandlerFunctionWithDelete(t *testing.T) {
	h1 := handler.NewMockHandler()
	chainOfHandlers := handler.NewChainOfHandlers(h1)
	fn := newKooperHandlerFunction(chainOfHandlers, "Deployment")

	err := fn(nil, &common.K8sEvent{
		Kind:      "Delete",
		HasSynced: true,
		Key:       "default/den-from-neverwhere",
		Object:    NewFakeDeployment(),
	})
	if err != nil {
		t.Error("err should be nil")
	}

	if h1.Called != true {
		t.Error("Handler should have been called")
	}
	if !reflect.DeepEqual(h1.PassedK8sManifest, []byte{}) {
		t.Error("K8sManifest have been set")
	}
	if !reflect.DeepEqual(h1.PassedPayload, []byte{}) {
		t.Error("Payload should be empty")
	}
}

func TestGetResourceFuncListFromConfig(t *testing.T) {
	configFile := path.Join(path.Dir(thisFilename), "fixtures", "config.toml")
	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--config=%s", configFile),
	}

	conf, err := config.NewConfig()
	if err != nil {
		t.Error("The config should have been parsed without errors")
	}
	resourceList, err := GetResourcesFuncListFromConfig(conf)
	if err != nil {
		t.Error("The resource list should have been returned without errors")
	}
	expected := 5
	if len(resourceList) != expected {
		t.Errorf("resourceList should have %d resource, have %d instead", expected, len(resourceList))
	}
}

func TestGetResourceWatcherList(t *testing.T) {
	resourcesFuncList := []func(ResourceWatcherArgs) watcher.ResourceWatcher{
		NewDeploymentWatcher,
		NewServiceWatcher,
		NewDaemonsetWatcher,
		NewIngressWatcher,
		NewStatefulsetWatcher,
	}
	rwl := GetResourceWatcherList(
		resourcesFuncList,
		ResourceWatcherArgs{
			Clientset:       nil,
			Namespace:       "",
			LabelSelector:   "",
			ChainOfHandlers: nil,
		})
	expected := 5
	if len(rwl) != expected {
		t.Errorf("resource watcher list should have %d resource, have %d instead", expected, len(rwl))
	}
}
