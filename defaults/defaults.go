package defaults

import _ "embed"

// DefaultPort is default for server to run on
// open ip tables using:
//  iptables -A INPUT -p tcp --dport 5050 -j ACCEPT
const DefaultPort = "5050"

// BigResponse is a large response, should take multiple data frames
// for the client to receive it. Total size is 275075 bytes
//go:embed bigresponse.txt
var BigResponse string
