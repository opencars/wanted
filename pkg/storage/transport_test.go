package storage

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Works properly.
func search(t Transport, id string) int {
	for i, x := range t {
		if x.ID == id {
			return i
		}
	}
	return -1
}

func TestTransport_Sort(t *testing.T) {
	transport := Transport{{ID: "9"}, {ID: "8"}, {ID: "7"}, {ID: "6"}, {ID: "5"}, {ID: "4"}, {ID: "3"}, {ID: "2"}, {ID: "1"}}
	expected := Transport{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}, {ID: "6"}, {ID: "7"}, {ID: "8"}, {ID: "9"}}

	sort.Sort(transport)
	assert.Equal(t, expected, transport)
}

func TestTransport_Search(t *testing.T) {
	rand.Seed(time.Now().Unix())

	arr := Transport(make([]WantedVehicle, 0, 10000))

	for i := 0; i < 10000; i++ {
		arr = append(arr, WantedVehicle{ID: strconv.Itoa(rand.Int())})
	}

	sort.Sort(arr)

	for i := 0; i < 10000; i++ {
		actual := arr.Search(arr[i].ID)
		expected := search(arr, arr[i].ID)
		if expected != actual {
			t.Errorf("failed, expected %d, got %d", expected, actual)
		}
	}
}

func BenchmarkTransport_Search(b *testing.B) {
	rand.Seed(time.Now().Unix())

	arr := Transport(make([]WantedVehicle, 0, b.N))

	for i := 0; i < b.N; i++ {
		arr = append(arr, WantedVehicle{ID: strconv.Itoa(rand.Int())})
	}

	sort.Sort(arr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		arr.Search(arr[i].ID)
	}
}

func BenchmarkTransport_search(b *testing.B) {
	rand.Seed(time.Now().Unix())

	arr := Transport(make([]WantedVehicle, 0, b.N))

	for i := 0; i < b.N; i++ {
		arr = append(arr, WantedVehicle{ID: strconv.Itoa(rand.Int())})
	}

	sort.Sort(arr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		search(arr, arr[i].ID)
	}
}
