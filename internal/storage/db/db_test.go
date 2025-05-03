package db

import (
	"os"
	"path/filepath"
	"testing"
	"yandex-go-advanced/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestGetRootDirectory(t *testing.T) {
	sgr := logger.Init()

	t.Run("should find root directory", func(t *testing.T) {
		tempDir := t.TempDir()
		goModPath := filepath.Join(tempDir, "go.mod")

		file, err := os.Create(goModPath)
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)

		currentDir, _ := os.Getwd()
		defer func() {
			err = os.Chdir(currentDir)
			assert.NoError(t, err)
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()
		err = os.Chdir(tempDir)
		assert.NoError(t, err)

		root, err := getRootDirectory()
		assert.NoError(t, err)
		assert.Contains(t, root, tempDir)
	})

	t.Run("should return error if go.mod not found", func(t *testing.T) {
		tempDir := t.TempDir()

		currentDir, _ := os.Getwd()
		defer func() {
			err := os.Chdir(currentDir)
			assert.NoError(t, err)
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()
		err := os.Chdir(tempDir)
		assert.NoError(t, err)

		root, err := getRootDirectory()
		assert.Error(t, err)
		assert.Empty(t, root)
		assert.Equal(t, "failed to find root directory", err.Error())
	})
}
