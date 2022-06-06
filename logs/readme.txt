Terminology and explanation:
remote: refers to a remote machine, running Linux (Centos7)
local: refers to a local machine, running MacOS

The IP address of the remote machine has been redacted, replaced with IP_REMOTE
The IP address of the local machine has been redacted, replaced with IP_LOCAL

To compile the client for the remote machine, go to the grpcbug dir and run:
% GOOS=linux go build -o clientbin client/main.go

To compile the server for the remote machine, go to the grpcbug dir and run:
% GOOS=linux go build -o serverbin server/main.go


This folder contains the following client-server logs:
local_client : client is on local machine, server is on remote. EOF observed.
local_client--work : client is on local machine, server is on remote. The -work=true option is passed to client. No EOF.
local_server : server is on local machine, client is on remote. No EOF.
local_server--work : server is on local machine, client is on remote. The -work=true option is passed to client. No EOF.
  (this was collected to serve as comparison against the local_client--work logs)

