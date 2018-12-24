package runtime

import "testing"

func TestNewClient(t *testing.T) {
	testEndpoint := "http://localhost:8080"
	c := NewClient(testEndpoint)
	if c == nil {
		t.Error("failed to init client")
	}
}
