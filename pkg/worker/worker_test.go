package worker_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/logger"
	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store/teststore"
	"github.com/opencars/wanted/pkg/worker"
)

func TestWorker_Parse_1(t *testing.T) {
	w := worker.New()

	reader, err := os.Open("../../testdata/14122019_2.json")
	assert.NoError(t, err)

	r, err := bom.NewReader(reader)
	assert.NoError(t, err)

	added, removed, err := w.Parse("14122019_2", r)
	assert.NoError(t, err)
	assert.Equal(t, 73551, len(added))
	assert.Equal(t, 0, len(removed))
}

func TestWorker_Parse_2(t *testing.T) {
	s := teststore.New()

	f, err := os.Open("../../testdata/14122019_1.json")
	assert.NoError(t, err)

	r, err := bom.NewReader(f)
	assert.NoError(t, err)

	vehicles := make([]model.Vehicle, 0)
	err = json.NewDecoder(r).Decode(&vehicles)
	assert.NoError(t, err)

	revision := model.TestRevision(t)
	revision.ID = "14122019_1"

	err = s.Vehicle().Create(revision, vehicles, nil)
	assert.NoError(t, err)

	w := worker.New()
	assert.NoError(t, w.Load(s))

	f, err = os.Open("../../testdata/14122019_2.json")
	assert.NoError(t, err)

	buff := &bytes.Buffer{}
	_, err = io.Copy(buff, f)
	assert.NoError(t, err)

	added1, removed1, err := w.Parse("14122019_2", bytes.NewReader(buff.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, 110, len(added1))
	assert.Equal(t, 163, len(removed1))

	added2, removed2, err := w.Parse("14122019_3", bytes.NewReader(buff.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(added2))
	assert.Equal(t, 0, len(removed2))
}

func TestMain(m *testing.M) {
	logger.Log = logger.NewLogger(ioutil.Discard)

	os.Exit(m.Run())
}
