package internal_test

import (
	"cloud.google.com/go/bigtable"
	"github.com/stretchr/testify/assert"

	"context"
	"testing"

	"github.com/pseudomuto/btsync/pkg/internal"
)

func TestWithBTServer(t *testing.T) {
	internal.WithBTServer(func(ac *bigtable.AdminClient) {
		tables, err := ac.Tables(context.Background())
		assert.Len(t, tables, 0)
		assert.NoError(t, err)
	})
}
