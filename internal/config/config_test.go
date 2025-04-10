package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Config initialized", func(t *testing.T) {
		cnf := Init()
		assert.NotNil(t, cnf)
	})
}

func TestInit_DefaultValues(t *testing.T) {
	t.Run("Config initialized by default", func(t *testing.T) {
		os.Clearenv()
		cnf := Init()

		assert.Equal(t, ":8080", *cnf.Server)
		assert.Equal(t, "http://localhost:8080", *cnf.BaseURL)
		assert.Equal(t, "", *cnf.FilePath)
		assert.Equal(t, "", *cnf.DataBase)
	})
}

func TestInit_ParseEnv(t *testing.T) {
	t.Run("Config initialized by env", func(t *testing.T) {
		err := os.Setenv("SERVER_ADDRESS", ":9090")
		assert.NoError(t, err)
		err = os.Setenv("BASE_URL", "http://localhost:9091")
		assert.NoError(t, err)
		err = os.Setenv("FILE_STORAGE_PATH", "/tmp/file")
		assert.NoError(t, err)
		err = os.Setenv("DATABASE_DSN", "postgres://localhost:5432/testdb")
		assert.NoError(t, err)

		cnf := Init()

		assert.Equal(t, ":9090", *cnf.Server)
		assert.Equal(t, "http://localhost:9091", *cnf.BaseURL)
		assert.Equal(t, "/tmp/file", *cnf.FilePath)
		assert.Equal(t, "postgres://localhost:5432/testdb", *cnf.DataBase)
	})
}
