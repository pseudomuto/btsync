package bt

import (
	"context"
)

// Table described a table's schema
type Table struct {
	Name     string
	Families map[string]GCPolicy
}

// A GCPolicy is the stringified representation of a column family's GC policy
type GCPolicy string

// A TableResult is a result type of table info
type TableResult struct {
	Table *Table
	Err   error
}

// Tables enumerates the tables in the BT cluster
func (c *Connection) Tables(ctx context.Context) <-chan TableResult {
	stream := make(chan TableResult)

	go func() {
		defer close(stream)

		tbls, err := c.ac.Tables(ctx)
		if err != nil {
			stream <- TableResult{Err: err}
			return
		}

		for _, t := range tbls {
			select {
			case <-ctx.Done():
				return
			case stream <- c.getTable(ctx, t):
			}
		}
	}()

	return stream
}

func (c *Connection) getTable(ctx context.Context, tbl string) TableResult {
	info, err := c.ac.TableInfo(ctx, tbl)
	if err != nil {
		return TableResult{Err: err}
	}

	table := &Table{Name: tbl, Families: make(map[string]GCPolicy)}
	for _, cf := range info.FamilyInfos {
		table.Families[cf.Name] = GCPolicy(cf.GCPolicy)
	}

	return TableResult{Table: table}
}
