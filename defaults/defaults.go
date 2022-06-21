package defaults

import (
	_ "embed"

	"google.golang.org/grpc"
)

// DefaultPort is default for server to run on
// open ip tables using:
//  iptables -A INPUT -p tcp --dport 5050 -j ACCEPT
const DefaultPort = "5050"

// BigResponse is a large response, should take multiple data frames
// for the client to receive it. Total size is 275075 bytes
//go:embed bigresponse.txt
var BigResponse string

// WorkaroundCfg returns server options that seem to prevent the client from
// experiencing unexpected EOFs
func WorkaroundCfg() []grpc.DialOption {
	// empirically determined to be enough overhead to avoid the EOF error
	// possibly related to the frame header being 9 bytes
	const windowOverhead = 9

	// default window size is 65535 from
	// google.golang.org/grpc/internal/transport/defaults.go
	// Exceeding this will disable dynamic window and bdp for the transport
	largerWindowSize := int32(len(BigResponse)) + windowOverhead
	return []grpc.DialOption{
		grpc.WithInitialWindowSize(largerWindowSize),
		grpc.WithInitialConnWindowSize(largerWindowSize),
	}
}
