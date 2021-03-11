Raft Leader Election Implementation
====================================

Running the Example
--------------------
To start the example, run the following command with the number of nodes you want in the test (Defaults to 7):
```bash
$ ./start.sh <number of nodes>
```

Alternatively, you can start individual nodes yourself:
```bash
$ go build
$ ./leader_election <number of nodes> <this node's ID>
```

> Note: Leader nodes automatically crash after sending a few heartbeats

Example Logs
-------------
Logs are stored in `output.txt`.
```
[21:48:46.4878 | ID 000 | F0 ]  Timed out waiting for heartbeat
[21:48:46.4883 | ID 000 | C1 ]  Requesting votes
[21:48:46.4889 | ID 002 | F0 ]  Timed out waiting for heartbeat
[21:48:46.4889 | ID 002 | C1 ]  Requesting votes
[21:48:46.4897 | ID 003 | F0 ]  Timed out waiting for heartbeat
[21:48:46.4898 | ID 003 | C1 ]  Requesting votes
[21:48:46.4902 | ID 001 | F0 ]  Received request from Node 0 for term 1
[21:48:46.4905 | ID 001 | F0 ]  Voting for Node 0 for term 1
[21:48:46.4915 | ID 001 | F1 ]  Received request from Node 2 for term 1
[21:48:46.4915 | ID 001 | F1 ]  Rejecting request from Node 2
[21:48:46.4918 | ID 004 | F0 ]  Received request from Node 0 for term 1
[21:48:46.4919 | ID 004 | F0 ]  Voting for Node 0 for term 1
[21:48:46.4919 | ID 000 | C1 ]  Received request from Node 2 for term 1
[21:48:46.4920 | ID 000 | C1 ]  Rejecting request from Node 2
[21:48:46.4925 | ID 004 | F1 ]  Received request from Node 2 for term 1
[21:48:46.4925 | ID 003 | C1 ]  Received request from Node 0 for term 1
[21:48:46.4926 | ID 004 | F1 ]  Rejecting request from Node 2
[21:48:46.4926 | ID 000 | C1 ]  Received request from Node 3 for term 1
[21:48:46.4927 | ID 000 | C1 ]  Rejecting request from Node 3
[21:48:46.4928 | ID 003 | C1 ]  Rejecting request from Node 0
[21:48:46.4929 | ID 000 | C1 ]  Received vote (2/5) from Node 1
[21:48:46.4930 | ID 000 | C1 ]  Received vote (3/5) from Node 4
[21:48:46.4930 | ID 000 | L1 ]  Is the Leader
[21:48:46.4930 | ID 000 | L1 ]  Sending heartbeats
[21:48:46.4930 | ID 003 | C1 ]  Received request from Node 2 for term 1
[21:48:46.4931 | ID 003 | C1 ]  Rejecting request from Node 2
[21:48:46.4931 | ID 002 | C1 ]  Received request from Node 0 for term 1
[21:48:46.4932 | ID 002 | C1 ]  Rejecting request from Node 0
[21:48:46.4933 | ID 001 | F1 ]  Received request from Node 3 for term 1
[21:48:46.4933 | ID 001 | F1 ]  Rejecting request from Node 3
[21:48:46.4936 | ID 004 | F1 ]  Received request from Node 3 for term 1
[21:48:46.4936 | ID 003 | C1 ]  Received No vote from Node 0
[21:48:46.4936 | ID 004 | F1 ]  Rejecting request from Node 3
[21:48:46.4937 | ID 003 | C1 ]  Received No vote from Node 1
[21:48:46.4937 | ID 002 | C1 ]  Received request from Node 3 for term 1
[21:48:46.4937 | ID 002 | C1 ]  Rejecting request from Node 3
[21:48:46.4938 | ID 003 | C1 ]  Received No vote from Node 4
[21:48:46.4939 | ID 001 | F1 ]  Received heartbeat
[21:48:46.4939 | ID 003 | C1 ]  Received No vote from Node 2
[21:48:46.4940 | ID 002 | C1 ]  Received No vote from Node 0
[21:48:46.4940 | ID 002 | C1 ]  Received No vote from Node 1
[21:48:46.4941 | ID 003 | F1 ]  Received heartbeat
[21:48:46.4943 | ID 002 | C1 ]  Received No vote from Node 3
[21:48:46.4943 | ID 004 | F1 ]  Received heartbeat
[21:48:46.5070 | ID 002 | C1 ]  Received No vote from Node 4
[21:48:46.5070 | ID 002 | F1 ]  Received heartbeat
[21:48:46.5946 | ID 000 | L1 ]  Sending heartbeats
[21:48:46.5948 | ID 001 | F1 ]  Received heartbeat
[21:48:46.5952 | ID 003 | F1 ]  Received heartbeat
[21:48:46.5954 | ID 004 | F1 ]  Received heartbeat
[21:48:46.5955 | ID 002 | F1 ]  Received heartbeat
[21:48:46.6957 | ID 000 | L1 ]  Sending heartbeats
[21:48:46.6959 | ID 001 | F1 ]  Received heartbeat
[21:48:46.6963 | ID 003 | F1 ]  Received heartbeat
[21:48:46.6964 | ID 004 | F1 ]  Received heartbeat
[21:48:46.6966 | ID 002 | F1 ]  Received heartbeat
...
```

These can be `grep`'ed for individual node logs with the following command, where `XXX` is the ID of the node front-padded with zeros (i.e., Node 0's ID would be `000`, Node 1's ID would be `001`, etc.):
```bash
$ grep "ID XXX" output.txt
```

