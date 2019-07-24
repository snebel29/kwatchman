package slack

import (
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMinInt(t *testing.T) {
	cases := []map[string]int{
		{"a": 1, "b": 2, "result": 1},
		{"a": 1, "b": -9, "result": -9},
	}
	for _, s := range cases {
		if minInt(s["a"], s["b"]) != s["result"] {
			t.Errorf("min of %d and %d != %d", s["a"], s["b"], s["result"])
		}
	}
}

func TestTruncate(t *testing.T) {
	type testCase struct {
		text     string
		limit    int
		expected string
	}

	tests := []testCase{
		{"abcdefgh", 5, "abcde"},
		{"abcdefgh", 1000, "abcdefgh"},
		{"我想玩电脑", 2, "我想"},
	}

	for _, test := range tests {
		result := truncateString(test.text, test.limit)
		if result != test.expected {
			t.Errorf("%s != %s", test.expected, result)
		}
	}
}

func TestSucessfullMsgToSlack(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{}"))
		if err != nil {
			t.Error(err)
		}
	}))
	defer testServer.Close()

	manifest := []byte("manifest")
	payload := []byte("payload")
	resourceKind := "Deployment"

	h := NewSlackHandler(config.Handler{
		ClusterName: "myClusterName",
		WebhookURL:  testServer.URL,
	})

	evt := &handler.Event{
		K8sEvt:       &common.K8sEvent{Kind: "Update"},
		RunNext:      true,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	}

	err := h.Run(nil, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.RunNext != true {
		t.Error("RunNext should be true")
	}
	if !reflect.DeepEqual(evt.Payload, payload) {
		t.Errorf("Payload %s should match %s", string(evt.Payload), string(payload))
	}
	if !reflect.DeepEqual(evt.K8sManifest, manifest) {
		t.Errorf("K8sManifest %s should match %s", string(evt.Payload), string(manifest))
	}
}

func TestFailedMsgToSlack(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte("{}"))
		if err != nil {
			t.Error(err)
		}
	}))
	defer testServer.Close()

	manifest := []byte("manifest")
	payload := []byte("payload")
	resourceKind := "Deployment"

	h := NewSlackHandler(config.Handler{
		ClusterName: "myClusterName",
		WebhookURL:  testServer.URL,
	})

	evt := &handler.Event{
		K8sEvt:       &common.K8sEvent{Kind: "Update"},
		RunNext:      true,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	}

	err := h.Run(nil, evt)
	if err == nil {
		t.Error(err)
	}

	if evt.RunNext != false {
		t.Error("RunNext should be false")
	}

}

func TestSlackHandler_buildTextField(t *testing.T) {
	text := buildTextField([]byte(""))
	if len(text) != 0 {
		t.Error("text should be empty")
	}
	message := "my difference text"
	text = buildTextField([]byte(message))
	expected := len(message) + 6
	if len(text) != expected {
		t.Errorf("text length should match got %d and %d instead", len(text), expected)
	}
}
