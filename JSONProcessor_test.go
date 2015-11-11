package negotiator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldProcessJSONAcceptHeader(t *testing.T) {
	var fibTests = []struct {
		acceptheader string // input
	}{
		{"application/json"},
		{"application/json-"},
		{"text/json"},
		{"+json"},
	}

	jsonProcessor := &JSONProcessor{}

	for _, tt := range fibTests {
		result := jsonProcessor.CanProcess(tt.acceptheader)
		assert.True(t, result, "Should process "+tt.acceptheader)
	}
}
