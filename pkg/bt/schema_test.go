package bt_test

import (
	"cloud.google.com/go/bigtable"
	"github.com/stretchr/testify/assert"

	"context"
	"testing"
	"time"

	"github.com/pseudomuto/btsync/pkg/bt"
	"github.com/pseudomuto/btsync/pkg/internal"
)

func TestTables(t *testing.T) {
	ctx := context.Background()

	t.Run("when connection is bad", func(t *testing.T) {
		conn, err := bt.New(ctx, bt.ConnectionConfig{
			Project:  "proj",
			Instance: "instance",
		})

		assert.NotNil(t, conn)
		assert.NoError(t, err)

		res := <-conn.Tables(ctx)
		assert.Nil(t, res.Table)
		assert.Error(t, res.Err)
	})

	t.Run("enumerates the tables in the cluster", func(t *testing.T) {
		internal.WithBTServer(func(ac *bigtable.AdminClient) {
			tables := []*bigtable.TableConf{
				{
					TableID: "no_column_families",
				},
				{
					TableID: "default_gc_policy",
					Families: map[string]bigtable.GCPolicy{
						"data": bigtable.NoGcPolicy(),
					},
				},
				{
					TableID: "mixed_gc_policies",
					Families: map[string]bigtable.GCPolicy{
						"maxVersions": bigtable.MaxVersionsPolicy(10),
						"maxAge":      bigtable.MaxAgePolicy(10 * time.Hour),
						"maxVersionAndAge": bigtable.IntersectionPolicy(
							bigtable.MaxVersionsPolicy(10),
							bigtable.MaxAgePolicy(10*time.Hour),
						),
						"maxVersionOrAge": bigtable.UnionPolicy(
							bigtable.MaxVersionsPolicy(10),
							bigtable.MaxAgePolicy(10*time.Hour),
						),
					},
				},
			}

			for _, table := range tables {
				assert.NoError(t, ac.CreateTableFromConf(ctx, table))
			}

			conn := bt.NewWithAdminClient(ac)
			defer conn.Close()

			cursor := conn.Tables(ctx)

			for _, table := range tables {
				res := <-cursor
				assert.Equal(t, table.TableID, res.Table.Name)
				assert.NoError(t, res.Err)

				for name, policy := range table.Families {
					assert.Equal(t, bt.GCPolicy(policy.String()), res.Table.Families[name])
				}
			}
		})
	})
}
