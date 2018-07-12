package config_test

import (
	"github.com/stretchr/testify/assert"

	"context"
	"testing"

	"github.com/pseudomuto/btsync/pkg/config"
)

func TestParseDirectory(t *testing.T) {
	t.Run("parses partition files within the supplied directory", func(t *testing.T) {
		tests := []struct {
			dir           string
			numPartitions int
		}{
			{"../../testdata/doesnt_exist", 0},
			{"../../testdata/partitions", 1},
		}

		ctx := context.Background()

		for _, test := range tests {
			count := 0

			for res := range config.ParseDirectory(ctx, test.dir) {
				assert.NotEmpty(t, res.Partition.Name)
				assert.NoError(t, res.Err)
				count++
			}

			assert.Equal(t, test.numPartitions, count)
		}
	})

	t.Run("returns an error for malformed files", func(t *testing.T) {
		tests := []struct {
			dir  string
			errs []string
		}{
			{
				"../../testdata/invalid",
				[]string{"yaml: unmarshal errors:"},
			},
		}

		for _, test := range tests {
			errIdx := 0
			errCount := 0

			for res := range config.ParseDirectory(context.Background(), test.dir) {
				assert.Empty(t, res.Partition)
				assert.Contains(t, res.Err.Error(), test.errs[errIdx])

				errIdx++
				errCount++
			}

			assert.Equal(t, len(test.errs), errCount)
		}
	})
}
