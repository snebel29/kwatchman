package slack

import (
	"fmt"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestCliArgs(t *testing.T) {
	clusterName := "myCluster"
	webhookURL := "http://myWebhookURL.com"
	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--cluster-name=%s", clusterName),
		fmt.Sprintf("--slack-webhook=%s", webhookURL),
	}

	kingpin.Parse()
	cli := newCLI()
	if cli.clusterName != clusterName {
		t.Errorf("%s != %s", cli.clusterName, clusterName)
	}
	if cli.webhookURL != webhookURL {
		t.Errorf("%s != %s", cli.webhookURL, webhookURL)
	}
}

func TestMsgToSlack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	s := ts.URL
	webhookURL = &s
	c := "myCluster"
	clusterName = &c

	evt := &common.K8sEvent{}
	manifest := []byte("manifest")
	payload := []byte("payload")
	resourceKind := "Deployment"

	h := NewSlackHandler()
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
