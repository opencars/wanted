package worker_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/worker"
)

func TestWorker_Parse_1(t *testing.T) {
	w := worker.New()

	reader, err := os.Open("../../testdata/14082020_1.json")
	assert.NoError(t, err)

	r, err := bom.NewReader(reader)
	assert.NoError(t, err)

	added, removed, err := w.Parse("14082020_1", r)
	assert.NoError(t, err)
	assert.Equal(t, 73689, len(added))
	assert.Equal(t, 0, len(removed))
}
