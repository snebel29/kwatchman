package slack

import (
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/config"
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

func TestMsgToSlack(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer testServer.Close()

	evt := &common.K8sEvent{}
	manifest := []byte("manifest")
	payload := []byte("payload")
	resourceKind := "Deployment"

	h := NewSlackHandler(config.Handler{
		ClusterName: "myClusterName",
		WebhookURL:  testServer.URL,
	})
	output, err := h.Run(nil, handler.Input{
		Evt:          evt,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	})
	if err != nil {
		t.Error(err)
	}

	if output.RunNext != true {
		t.Error("RunNext should be true")
	}
	if !reflect.DeepEqual(output.Payload, payload) {
		t.Errorf("Payload %s should match %s", string(output.Payload), string(payload))
	}
	if !reflect.DeepEqual(output.K8sManifest, manifest) {
		t.Errorf("K8sManifest %s should match %s", string(output.Payload), string(manifest))
	}
}
