package related

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/brianneville/grpcbug/defaults"
	"github.com/brianneville/grpcbug/proto"
	"github.com/stretchr/testify/require"
)

const (
	responseTiny    = iota // response that is at least < 65535 bytes
	responseBig            // response that requires multiple messages
	responseExhaust        // response that exceeds the default grpc.max_receive_message_length
)

// default max size of message that can be received
const maxMsgRecvSize = 4194304

// mockServer mocks a MockRPC Server.
type mockServer struct {
	addr string
	grpc *grpc.Server
	*proto.UnimplementedMockRPCServer

	responseSize int
}

func (m *mockServer) Stop()        { m.grpc.Stop() }
func (m *mockServer) Addr() string { return m.addr }

func (m *mockServer) WithResponse(responseSize int, f func()) {
	defer func(s int) { m.responseSize = s }(m.responseSize)
	m.responseSize = responseSize
	f()
}

func (m mockServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {

	var contents []byte
	switch m.responseSize {
	case responseTiny:
		contents = []byte(":^)")
	case responseBig:
		contents = []byte(defaults.BigResponse)
	case responseExhaust:
		var arr [maxMsgRecvSize]byte
		contents = arr[:]
	}

	return &proto.GetResponse{
		B: contents,
	}, nil
}

func newMockServer(addr string) *mockServer {
	m := &mockServer{
		addr:                       addr,
		grpc:                       grpc.NewServer(),
		UnimplementedMockRPCServer: &proto.UnimplementedMockRPCServer{},
	}
	proto.RegisterMockRPCServer(m.grpc, m)
	return m
}

func startMockServer(t *testing.T) *mockServer {
	l, err := net.Listen("tcp", "")
	require.NoError(t, err)
	return startMockServerWithListener(t, l)
}

// startMockServerWithListener starts a new mock server with the provided listener
func startMockServerWithListener(t *testing.T, l net.Listener) *mockServer {
	srv := newMockServer(l.Addr().String())
	t.Cleanup(srv.grpc.Stop)

	go func() {
		require.NoError(t, srv.grpc.Serve(l))
	}()
	return srv
}

func (m *mockServer) NewClient(ctx context.Context, work bool) (proto.MockRPCClient, error) {

	// dial options applied in TestListResources from github.com/gravitational/teleport/api/client
	// 1. insecure credentials
	// 2. teleport test sets WithKeepaliveParams -- this does not affect the error so is omitted
	// 3. teleport test sets WithBlock -- this does not affect the error but include it anyway
	// since we should wait until the connection is up before doing anything.

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	if work {
		dialOpts = append(dialOpts, defaults.WorkaroundCfg()...)
	}

	conn, err := grpc.DialContext(ctx, m.Addr(), dialOpts...)
	if err != nil {
		return nil, err
	}
	return proto.NewMockRPCClient(conn), nil
}
func requireExhausted(t *testing.T, err error) {
	if e, ok := status.FromError(err); ok {
		if e.Code() == codes.ResourceExhausted {
			return
		}
	}
	t.Fatal(err)
}

// TestExhaustConn tests the behaviour of a grpc client/server immediately after
// the client has received a response which exceeded its grpc.max_receive_message_length
// (default 4194304 bytes). Minimal reproduction of gravitational/teleport#13548
/*
results:

big_response testcase sees EOF reliably in 5K runs:
$ go test github.com/brianneville/grpcbug/related \
 -run "^\QTestExhaustConn\E$/^\Qbig_response\E$" -count=5000 --failfast

other two testcases did not see any EOF even after 100K runs:
*/
func TestExhaustConn(t *testing.T) {
	for _, tc := range []struct {
		name           string
		secondResponse int
		useWorkaround  bool
	}{{
		name:           "tiny_response",
		secondResponse: responseTiny,
	}, {
		name:           "big_response",
		secondResponse: responseBig,
	}, {
		name:           "big_response_work",
		secondResponse: responseBig,
		useWorkaround:  true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			srv := startMockServer(t)
			clt, err := srv.NewClient(ctx, tc.useWorkaround)
			require.NoError(t, err)

			srv.WithResponse(responseExhaust, func() {
				_, err = clt.Get(ctx, &proto.GetRequest{})
				require.Error(t, err)
				requireExhausted(t, err)
			})

			srv.WithResponse(tc.secondResponse, func() {
				_, err = clt.Get(ctx, &proto.GetRequest{})
				require.NoError(t, err)
			})
		})
	}
}
