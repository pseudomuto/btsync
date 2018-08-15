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

// Table fetches details about the specified table
func (c *connection) Table(ctx context.Context, name string) (*Table, error) {
	info, err := c.ac.TableInfo(ctx, name)
	if err != nil {
		return nil, err
	}

	table := &Table{Name: name, Families: make(map[string]GCPolicy)}
	for _, cf := range info.FamilyInfos {
		table.Families[cf.Name] = GCPolicy(cf.GCPolicy)
	}

	return table, nil
}

// Tables enumerates the tables in the BT cluster
func (c *connection) Tables(ctx context.Context) <-chan TableResult {
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

func (c *connection) getTable(ctx context.Context, tbl string) TableResult {
	table, err := c.Table(ctx, tbl)
	return TableResult{Table: table, Err: err}
}
