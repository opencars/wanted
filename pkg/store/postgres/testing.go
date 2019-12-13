package postgres

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/config"
)

func TestDB(t *testing.T, settings *config.Settings) (*Store, func(...string)) {
	t.Helper()

	store, err := New(settings)
	assert.NoError(t, err)

	return store, func(tables ...string) {
		if len(tables) > 0 {
			_, err = store.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			assert.NoError(t, err)
		}

		store.db.Close()
	}
}
