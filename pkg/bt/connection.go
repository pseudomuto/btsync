package bt

import (
	"cloud.google.com/go/bigtable"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"context"
)

// ConnectionConfig represents connection details for BigTable clusters
type ConnectionConfig struct {
	Project  string
	Instance string
	JWT      []byte
	Options  []option.ClientOption
}

// A Connection represents an admin connection to a BigTable cluster
type Connection struct {
	ac *bigtable.AdminClient
}

// New returns a new connection to a BigTable cluster
func New(ctx context.Context, cfg ConnectionConfig) (*Connection, error) {
	opts := make([]option.ClientOption, 0, len(cfg.Options)+1)
	opts = append(opts, cfg.Options...)

	if cfg.JWT != nil {
		token, err := google.JWTConfigFromJSON(cfg.JWT, bigtable.Scope)
		if err != nil {
			return nil, err
		}

		opts = append(opts, option.WithTokenSource(token.TokenSource(ctx)))
	}

	ac, err := bigtable.NewAdminClient(ctx, cfg.Project, cfg.Instance, opts...)
	if err != nil {
		return nil, err
	}

	return NewWithAdminClient(ac), nil
}

// NewWithAdminClient returns a connection that uses the supplied admin client for the underlying connection
func NewWithAdminClient(ac *bigtable.AdminClient) *Connection {
	return &Connection{ac: ac}
}

// Close closes the underlying BT client
func (c *Connection) Close() error {
	return c.ac.Close()
}
