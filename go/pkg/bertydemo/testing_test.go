package bertydemo

import (
	context "context"
	"testing"

	"berty.tech/berty/go/internal/grpcutil"
	"berty.tech/berty/go/internal/ipfsutil"

	grpc "google.golang.org/grpc"
)

type cleanFunc func()

func testingInMemoryClient(t *testing.T) (*Client, ipfsutil.CoreAPIMock, cleanFunc) {
	t.Helper()

	ctx := context.Background()
	ipfsmock := ipfsutil.TestingCoreAPI(ctx, t)
	opts := &Opts{
		CoreAPI:          ipfsmock,
		OrbitDBDirectory: ":memory:",
	}

	demo, err := New(opts)
	if err != nil {
		t.Fatal(err)
	}

	return demo, ipfsmock, func() {
		demo.Close()
		ipfsmock.Close()
	}
}

func testingClientService(t *testing.T, srv DemoServiceServer) (DemoServiceClient, cleanFunc) {
	t.Helper()

	listener := grpcutil.NewPipeListener()

	server := grpc.NewServer()
	RegisterDemoServiceServer(server, srv)

	conn, err := listener.NewClientConn(grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	go server.Serve(listener)

	return NewDemoServiceClient(conn), func() {
		listener.Close()
	}
}
