package store

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/model"
)

// Works properly.
func search(t Transport, id string) int {
	for i, x := range t {
		if x.CheckSum == id {
			return i
		}
	}
	return -1
}

func TestTransport_Sort(t *testing.T) {
	transport := Transport{{CheckSum: "9"}, {CheckSum: "8"}, {CheckSum: "7"}, {CheckSum: "6"}, {CheckSum: "5"}, {CheckSum: "4"}, {CheckSum: "3"}, {CheckSum: "2"}, {CheckSum: "1"}}
	expected := Transport{{CheckSum: "1"}, {CheckSum: "2"}, {CheckSum: "3"}, {CheckSum: "4"}, {CheckSum: "5"}, {CheckSum: "6"}, {CheckSum: "7"}, {CheckSum: "8"}, {CheckSum: "9"}}

	sort.Sort(transport)
	assert.Equal(t, expected, transport)
}

func TestTransport_Search(t *testing.T) {
	rand.Seed(time.Now().Unix())

	arr := Transport(make([]model.Vehicle, 0, 10000))

	for i := 0; i < 10000; i++ {
		arr = append(arr, model.Vehicle{CheckSum: strconv.Itoa(rand.Int())})
	}

	sort.Sort(arr)

	for i := 0; i < 10000; i++ {
		actual := arr.Search(arr[i].CheckSum)
		expected := search(arr, arr[i].CheckSum)
		if expected != actual {
			t.Errorf("failed, expected %d, got %d", expected, actual)
		}
	}
}

func BenchmarkTransport_Search(b *testing.B) {
	rand.Seed(time.Now().Unix())

	arr := Transport(make([]model.Vehicle, 0, b.N))

	for i := 0; i < b.N; i++ {
		arr = append(arr, model.Vehicle{CheckSum: strconv.Itoa(rand.Int())})
	}

	sort.Sort(arr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		arr.Search(arr[i].CheckSum)
	}
}

func BenchmarkTransport_search(b *testing.B) {
	rand.Seed(time.Now().Unix())

	arr := Transport(make([]model.Vehicle, 0, b.N))

	for i := 0; i < b.N; i++ {
		arr = append(arr, model.Vehicle{CheckSum: strconv.Itoa(rand.Int())})
	}

	sort.Sort(arr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		search(arr, arr[i].CheckSum)
	}
}
