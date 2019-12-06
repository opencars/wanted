package postgres

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T, settings *config.Database) (*Store, func(...string)) {
	t.Helper()

	store, err := New(
		settings.Host, settings.Port,
		settings.User, settings.Password,
		settings.Name,
	)

	require.NoError(t, err)
	return store, func(tables ...string) {
		if len(tables) > 0 {
			_, err = store.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			assert.NoError(t, err)
		}

		store.db.Close()
	}
}
