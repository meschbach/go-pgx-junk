package go_pgx_junk

import (
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"testing/fstest"
)

func TestLoadPGFromFS(t *testing.T) {
	t.Run("parses pg-connection as DSN", func(t *testing.T) {
		exampleFS := fstest.MapFS{
			"pg-connection": &fstest.MapFile{
				Data: []byte("port=5482 host=example.pg.local"),
			},
		}

		result, err := LoadFromFS(exampleFS)
		require.NoError(t, err)

		assert.Equal(t, "example.pg.local", result.Host)
		assert.Equal(t, uint16(5482), result.Port)
	})

	t.Run("parses pg-connection as URL", func(t *testing.T) {
		exampleFS := fstest.MapFS{
			"pg-connection": &fstest.MapFile{
				Data: []byte("postgres://neutron.scale:9876"),
			},
		}

		result, err := LoadFromFS(exampleFS)
		require.NoError(t, err)

		assert.Equal(t, "neutron.scale", result.Host)
		assert.Equal(t, uint16(9876), result.Port)
	})

	t.Run("loads from individual files", func(t *testing.T) {
		dbName := faker.Name()
		parsedURL, err := url.Parse(faker.URL())
		require.NoError(t, err)
		hostName := parsedURL.Hostname()

		exampleFS := fstest.MapFS{
			"pg-host": &fstest.MapFile{
				Data: []byte(hostName),
			},
			"pg-port": &fstest.MapFile{
				Data: []byte("65535"),
			},
			"pg-database": &fstest.MapFile{
				Data: []byte(dbName),
			},
		}

		result, err := LoadFromFS(exampleFS)
		require.NoError(t, err)

		assert.Equal(t, hostName, result.Host)
		assert.Equal(t, uint16(65535), result.Port)
		assert.Equal(t, dbName, result.Database)
	})
}
