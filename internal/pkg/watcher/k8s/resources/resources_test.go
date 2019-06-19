package resources

import (
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/handler/log"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

type fakeDeployment struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

func k8sIndividualResourceWatcherHelper(w watcher.ResourceWatcher, t *testing.T) {
	rw := w.(*K8sResourceWatcher)

	if rw.kind == "" {
		t.Errorf("kind should be != than \"\"")
	}
	if rw.ctrl == nil {
		t.Errorf("ctrl should be != nil")
	}
	if rw.stopC == nil {
		t.Errorf("stopC should be != nil")
	}
}

func TestK8sDeploymentWatcher(t *testing.T) {
	chainOfHandlers := handler.NewChainOfHandlers(log.NewLogHandler())
	k8sIndividualResourceWatcherHelper(NewK8sDeploymentWatcher(nil, "", chainOfHandlers), t)
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

	r, err = getManifest(&appsv1.DaemonSet{})
	if err == nil {
		t.Error("err should be error")
	}
}

func TestNewKooperHandlerFunctionWithAdd(t *testing.T) {
	h1 := handler.NewMockHandler()
	chainOfHandlers := handler.NewChainOfHandlers(h1)
	fn := NewKooperHandlerFunction(chainOfHandlers, "Deployment")

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
	fn := NewKooperHandlerFunction(chainOfHandlers, "Deployment")

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
