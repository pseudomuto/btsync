package cmd_test

import (
	"github.com/stretchr/testify/assert"

	"bytes"
	"strings"
	"testing"

	"github.com/pseudomuto/btsync/cmd/btsync/cmd"
)

func TestNewTableWriter(t *testing.T) {
	t.Run("writes formatted table", func(t *testing.T) {
		buf := new(bytes.Buffer)

		table := cmd.NewTableWriter(buf)
		table.AppendBulk([][]string{
			{"Item1 Name", "Item1 Value"},
			{"Item2 Name", "Item2 Value"},
		})

		table.SetHeader([]string{"Name", "Value"})
		table.Render()

		expected := `
		|    NAME    |    VALUE    |
		|------------|-------------|
		| Item1 Name | Item1 Value |
		| Item2 Name | Item2 Value |
		`

		expected = strings.Replace(expected, "\t", "", -1)
		expected = strings.Replace(expected, "\n", "", 1)

		assert.Equal(t, expected, buf.String())
	})
}
