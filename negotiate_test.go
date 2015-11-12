package negotiator

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldAppendCustomResponseProcessors(t *testing.T) {

	var fakeResponseProcessor = &fakeProcessor{}

	New(fakeResponseProcessor)

	assert.Equal(t, 3, len(processors))
}

type fakeProcessor struct {
}

func (*fakeProcessor) CanProcess(mediaRange string) bool {
	return true
}

func (*fakeProcessor) Process(w http.ResponseWriter, model interface{}) {

}
