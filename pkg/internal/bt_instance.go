package internal

import (
	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/bigtable/bttest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"context"
)

// WithBTServer creates an in-memory BT server and starts the server on a random port on the local machine. This is
// useful for testing
func WithBTServer(testFunc func(*bigtable.AdminClient)) {
	svr, _ := bttest.NewServer("localhost:0")
	conn, _ := grpc.Dial(svr.Addr, grpc.WithInsecure())

	ac, _ := bigtable.NewAdminClient(context.Background(), "proj", "instance", option.WithGRPCConn(conn))
	defer ac.Close()

	testFunc(ac)
}
