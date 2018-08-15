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

func TestTable(t *testing.T) {
	ctx := context.Background()

	t.Run("when table exists", func(t *testing.T) {
		tbl := &bigtable.TableConf{
			TableID: "some_table",
			Families: map[string]bigtable.GCPolicy{
				"data": bigtable.NoGcPolicy(),
			},
		}

		internal.WithBTServer(func(ac *bigtable.AdminClient) {
			// create the table
			assert.NoError(t, ac.CreateTableFromConf(ctx, tbl))

			conn := bt.NewWithAdminClient(ac)
			defer conn.Close()

			table, err := conn.Table(ctx, tbl.TableID)
			assert.NoError(t, err)
			assert.Len(t, table.Families, 1)
			assert.Equal(t, bt.GCPolicy(""), table.Families["data"])
		})
	})

	t.Run("when table not found", func(t *testing.T) {
		internal.WithBTServer(func(ac *bigtable.AdminClient) {
			conn := bt.NewWithAdminClient(ac)
			defer conn.Close()

			table, err := conn.Table(ctx, "some_table")
			assert.Nil(t, table)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "NotFound")
		})
	})
}

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
