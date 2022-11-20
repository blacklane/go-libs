package events

import (
	"strings"
	"testing"
	"time"
)

func FuzzWithBootstrapServers(f *testing.F) {
	f.Add("server1", "server2", "server3")
	f.Add("s1", "s2", "s3")
	f.Fuzz(func(t *testing.T, a string, b string, c string) {
		servers := []string{a, b, c}
		cfg := NewKafkaConfig(
			WithBootstrapServers(servers),
		)
		got := (*cfg)["bootstrap.servers"]
		want := strings.Join(servers, ",")
		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	})
}

func TestWithSessionTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    int
	}{
		{
			name:    "1-second",
			timeout: 1 * time.Second,
			want:    1 * 1000,
		},
		{
			name:    "1-minute",
			timeout: 1 * time.Minute,
			want:    1 * 60 * 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewKafkaConfig(
				WithSessionTimeout(tt.timeout),
			)
			if got := (*cfg)["session.timeout.ms"]; got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestWithKeyValue(t *testing.T) {
	tests := []struct {
		key   string
		value interface{}
	}{
		{
			key:   "int",
			value: int(12),
		},
		{
			key:   "string",
			value: "value",
		},
		{
			key:   "bool",
			value: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			cfg := NewKafkaConfig(
				WithKeyValue(tt.key, tt.value),
			)
			if got := (*cfg)[tt.key]; got != tt.value {
				t.Errorf("got: %v, want: %v", got, tt.value)
			}
		})
	}
}
