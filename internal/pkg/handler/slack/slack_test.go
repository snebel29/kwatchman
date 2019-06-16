package slack

import (
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
	//TODO: We need to mock nslopes/slack
}
