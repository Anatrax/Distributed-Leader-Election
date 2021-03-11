package main

import (
	"os"
	"log"
	"fmt"
    "time"
	"strconv"
	"net"
	"net/rpc"
	"net/http"
    "math/rand"
)

const BASE_ADDR int = 8800 // Start of port range for nodes to listen on
var NUM_NODES int          // Total number of nodes in the cluster

type MessageType int
const (
    RequestVote MessageType = iota
    YesVote
    NoVote
    AppendEntries
)

type Message struct {
    Type MessageType // The type of message
    SenderID int     // The ID of the node that sent the message
    Term int         // The current term of the node that sent the message
}

type Node struct {
    id int                  // Node ID
    channels []*rpc.Client  // Communication channels between other nodes
    msg_buffer chan Message // Buffer for incoming message
    state string            // Node state (Follower, Candidate, Leader)
    term int                // Node's current voting term
}

// Custom logging function 
func (node *Node) log(str string){
	time := time.Now()

    var state_char string
    switch this.state {
    case "Follower":
        state_char = "F"
    case "Candidate":
        state_char = "C"
    case "Leader":
        state_char = "L"
    }

    fmt.Printf(
        "[%s | ID %03d | %s%d ]  %s\n",
        time.Format("15:04:05.0000"),
        node.id,
        state_char,
        node.term,
        str,
    )
}

// Node constructor
func initNode() *Node {
    // Read in arguments
    node_id := parseArgs()

    // Set up node state
	node := new(Node)
	node.id = node_id
    node.state = "Follower" // "Nodes start as followers."
    node.term = 0
    node.msg_buffer = make(chan Message, 1000)

    // Set up communications
	node.initServer()
    node.initChannels()

	return node
}

// Get node's address
func (node *Node) getAddr(id int) string {
    return ":"+strconv.Itoa(BASE_ADDR + id)
}

// Start listening for connections from other nodes
func (node *Node) initServer(){
    // Register handlers
	rpc.Register(node)
	rpc.HandleHTTP()

    // Start listing at (arbitrary) BASE_PORT
	l, err := net.Listen("tcp", node.getAddr(node.id))
	if err != nil {
		log.Fatal("could not initialize server", err)
	}
	go http.Serve(l, nil)
}

// Create connections to each other node in the network,
// and store their information in the clients array
func (node *Node) initChannels(){
    for i := 0; i < NUM_NODES; i++ {
        // Attempt to initiate connections with each node
        client, err := rpc.DialHTTP("tcp", node.getAddr(i))
		if err != nil {
            // Retry indefinitely if connection fails
			for err!=nil{
                client, err = rpc.DialHTTP("tcp", node.getAddr(i))
			}
		}
        node.channels = append(node.channels, client)
	}
	return
}

// Multicast synchronous RPC
func (node *Node) multicast(msg_type MessageType, msg int) string {
    var res string
    for i := 0; i < NUM_NODES; i++ {
        if i != node.id {
            res = node.send(msg_type, msg, i)
        }
    }
    return res
}

// Send synchronous RPC
func (node *Node) send(msg_type MessageType, msg int, id int) string {
//    mtype_strings := [4]string{"RequestVote", "Yes vote", "No vote", "Heartbeat"}
//    node.log(fmt.Sprintf("Sending %s to Node %d", mtype_strings[msg_type], id))

    req := Message{Type:msg_type, SenderID:node.id, Term:msg}
    var res string
    err := node.channels[id].Call("Node.Post", req, &res)
    if err != nil {
        fmt.Println("Unable to send ping\n", err);
    }
    return res
}

// Add incoming message to buffer
func (node *Node) Post(req Message, res *string) error {
    node.msg_buffer <- req
    *res = node.getAddr(node.id)
	return nil
}

// Get command line arguments 
func parseArgs() int {
    // Use first argument to initialize global NUM_NODES variable
    num_nodes, err := strconv.Atoi(os.Args[1])
	if err!=nil {
		log.Fatal("error parsing command line arguments", err)
	}
    NUM_NODES = num_nodes

    // Use second argument to get this node's ID
    node_id, err := strconv.Atoi(os.Args[2])
	if err!=nil {
		log.Fatal("error parsing command line arguments", err)
	}
    return node_id
}

var this *Node // This node

func main(){
    this = initNode() // Initialize this node

    vote_count := 0
    timeout_counter := 0 // For ignoring cancelled timeouts

    for {
        switch this.state {
        case "Follower":
            //this.log("Waiting for heartbeat")
        case "Candidate":
            // "[The candidate] votes for itself..."
            vote_count = 1

            // "...and issues RequestVotes RPC's...to the other nodes in the cluster."
            this.log("Requesting votes")
            this.multicast(RequestVote, this.term)
        case "Leader":
            this.log("Is the Leader")

            // TEST: Crash the leader after a few heartbeats
            go func() {
                time.Sleep(1*time.Second)
                this.log("TEST: Crashing Leader Node")
                this.state = "Crashed"
            }()

            // "Leaders send periodic heartbeats
            // (AppendEntries RPC's that carry no log entries)"
            for ; this.state == "Leader"; {
                this.log("Sending heartbeats")
                this.multicast(AppendEntries, this.term)
                time.Sleep(100*time.Millisecond)
            }

            // TEST: Crashed state
            for {
            }
        }

        // Start timer
        election_timeout := make(chan int, 100)
        setTimeout(&election_timeout, &timeout_counter)

        for transition_flag := false; transition_flag == false; {
            select {
            case refs := <-election_timeout:
                if refs == 0 {
                    switch this.state {
                    case "Follower":
                        // "If a follower receives no communication during the election timeout,
                        // then it assumes there is no viable leader and begins an election
                        // to choose a leader"
                        this.log("Timed out waiting for heartbeat")
                        // "To begin an election, a follower increments its current term..."
                        this.term++
                        // "...and transitions to candidate state."
                        this.state = "Candidate"
                    case "Candidate":
                        // "If many followers become candidates at the same time, votes
                        // could be split so that no candidate obtains a majority."
                        this.log("Timed out waiting for votes")

                        // "When this happens, each candidate will time out and start a
                        // new election by incrementing its term..."
                        this.term++
                        // "and [loops back around to] initiating a new round
                        // of RequestVote RPC's"
                    }
                    transition_flag = true
                }
            case msg := <-this.msg_buffer:
                switch msg.Type {
                case RequestVote:
                    // "Each node will vote for at most one candidate in a given term,
                    // on a first-come-first-serve basis."
                    // TODO: If doing logger replication, last log entry term must be
                    //       greater than this node's last log entry term, or else the
                    //       length of the logs must be greater than this node's logs.
                    this.log(fmt.Sprintf(
                        "Received request from Node %d for term %d",
                        msg.SenderID,
                        msg.Term,
                    ))
                    if msg.Term > this.term {
                        // "Nodes remain followers as long as they receive valid RPC's
                        // from a leader or candidate"
                        transition_flag = true

                        this.log(fmt.Sprintf(
                            "Voting for Node %d for term %d",
                            msg.SenderID,
                            msg.Term,
                        ))
                        this.send(YesVote, this.term, msg.SenderID)
                        this.state = "Follower"
                        this.term = msg.Term
                    } else {
                        this.log(fmt.Sprintf(
                            "Rejecting request from Node %d",
                            msg.SenderID,
                        ))
                        this.send(NoVote, this.term, msg.SenderID)
                    }
                case YesVote:
                    vote_count++
                    this.log(fmt.Sprintf(
                        "Received vote (%d/%d) from Node %d",
                        vote_count,
                        NUM_NODES,
                        msg.SenderID,
                    ))

                    // "A candidate wins an election if it receives votes from
                    // a majority of the nodes in the cluster for the same term."
                    if vote_count >= NUM_NODES - NUM_NODES/2 {
                        // "Once a candidate wins an election, it becomes a leader."
                        this.state = "Leader"
                        transition_flag = true
                    }
                case NoVote:
                    this.log(fmt.Sprintf(
                        "Received No vote from Node %d",
                        msg.SenderID,
                    ))
                case AppendEntries:
                    switch this.state {
                    case "Follower":
                        // "Nodes remain followers as long as they receive valid RPC's
                        // from a leader or candidate"
                        transition_flag = true
                        // TODO: Append entries, update term to match
                    case "Candidate":
                        // "While waiting for votes, a candidate may receive an
                        // AppendEntries RPC's from another node claiming to be leader."
                        if msg.Term >= this.term {
                            // "If the leader's term (included in RPC) is at least
                            // as large as the candidate's current term,
                            // then the candidate recognizes the leader as legitimate
                            // and returns to follower state."

                            this.term = msg.Term
                            this.state = "Follower"
                        }
                        // "If the term in the RPC is smaller than the candidate's
                        // current term, then the candidate rejects the RPC
                        // and continues in candidate state."
                    }
                    this.log("Received heartbeat")//this.log("Received entries to append")
                }
            default: // Non-blocking select
            }
            if timeout_counter < 0 {
                this.log("ERROR: Negative timeout references, halting Node")
                for {
                }
            }
        }
    }
}

// Starts timeout timer, decrements counter and sends it as the timeout message
func setTimeout(timeout *chan int, timeout_counter *int) {
    *timeout_counter++
    go func() {
        // "Raft uses randomized election timeouts to ensure that split votes are rare
        // and that they are resolved quickly...election timeouts are chosen randomly
        // from a fixed interval (e.g., 150-300ms)."
        duration := 150 + rand.Intn(150)
        time.Sleep(time.Duration(duration) * time.Millisecond)
        *timeout_counter--
        *timeout <-*timeout_counter
    }()
}

