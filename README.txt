This repo reproduces a bug in gRPC

This bug is reproduced for the following setup:
 - a laptop running MacOS, sw_ver is:
  ProductName: macOS
  ProductVersion: 11.6
  BuildVersion: 20G165

 - a remote peer running Centos 7, where this server is on a high ping
   (~150ms should be enough)

 - modules and go version as seen in go.mod/go.sum

1. on laptop compile server/main.go:
 % GOOS=linux go build -o serverbin server/main.go

2. scp serverbin to the remote and start it.
 This will start a grpc server running on port 5050.
 the iptables may need to be updated for this port to be opened.

3. run the client:
 % go run client/main.go -addr=$REMOTE
The client will then return the following error:
  rpc error: code = Internal desc = unexpected EOF

Note that the client-side workaround (disable BDP estimation and dynamic window)
can be seen by running:
 % go run client/main.go -addr=$REMOTE -work=true