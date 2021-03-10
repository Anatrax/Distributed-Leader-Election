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
//    "math/rand"
)

const BASE_ADDR int = 9000 // Start of port range for nodes to listen on
var NUM_NODES int // Total number of nodes in the cluster

type MessageType int
const (
    RequestVote MessageType = iota
    YesVote
    NoVote
    AppendEntries
)

type Message struct {
    Type MessageType
    Sender string
    Value int
}

// We can add variables to this struct for node state
type Node struct {
    id int // Node ID
    channels []*rpc.Client // Communication channels between other nodes
    msg_buffer chan Message // Buffer for incoming message
    term int // Voting term
    state string // Node state (Follower, Candidate, Leader)
    hasVotedThisTerm bool
    vote_count int
}

// Custom node logging function 
func (node *Node) log(str string){
	time := time.Now()
    // Print the time, node id, and provided message
	fmt.Printf("[%s | ID %03d]  %s\n", time.Format("15:04:05.0000"), node.id, str)
}

// Node constructor
func init_node() *Node {
    node_id := parseArgs()
    fmt.Println("Initializing Node", node_id)

	node := new(Node)
	node.id = node_id
    node.state = "Follower"
    node.term = 0
    node.hasVotedThisTerm = false
    node.msg_buffer = make(chan Message, 1000)

	node.init_server()
    node.init_channels()
	node.log("initialized server and connected to all channels")

	return node
}

// Start listening for connections from other nodes
func (node *Node) init_server(){
    // Register handlers
    // RPC register will export all functions with of type "Node" that have a captial first letter
    // e.g. Pong will be exported, ping will not
	rpc.Register(node)
	rpc.HandleHTTP()

    // Start listing at (arbitrary) BASE_PORT
	l, err := net.Listen("tcp", id_to_addr(node.id))
	if err != nil {
		log.Fatal("could not initialize server", err)
	}
	go http.Serve(l, nil)
}

// Create connections to each other node in the network,
// and store their information in the clients array
func (node *Node) init_channels(){
    for i := 0; i < NUM_NODES; i++ {
        // Attempt to initiate connections with each node
        client, err := rpc.DialHTTP("tcp", id_to_addr(i))
		if err != nil {
            // Retry indefinitely if connection fails
			for err!=nil{
                client, err = rpc.DialHTTP("tcp", id_to_addr(i))
			}
		}
        node.log(fmt.Sprintf("opened channel to node %d", i))
        node.channels = append(node.channels, client)
	}
	return
}

// TODO: Rename broadcast to multicast
// Multicast synchronous RPC
func (node *Node) broadcast(msg_type MessageType, msg int) string {
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
    mtype_strings := [4]string{"RequestVote", "Yes vote", "No vote", "Heartbeat"}
    node.log(fmt.Sprintf("Sending %s to node %d", mtype_strings[msg_type], id))

    req := Message{Type:msg_type, Sender:id_to_addr(node.id), Value:msg}
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
    *res = id_to_addr(node.id)
	return nil
}

func id_to_addr(id int) string {
    return ":"+strconv.Itoa(BASE_ADDR + id)
}

// Get command line arguments 
func parseArgs() int {
    num_nodes, err := strconv.Atoi(os.Args[1])
	if err!=nil {
		log.Fatal("error parsing command line arguments", err)
	}
    NUM_NODES = num_nodes

    node_id, err := strconv.Atoi(os.Args[2])
	if err!=nil {
		log.Fatal("error parsing command line arguments", err)
	}
    return node_id
}

var this *Node
func main(){
    this = init_node()
    this.broadcast(RequestVote, this.id)
    select {
    case msg := <-this.msg_buffer:
        fmt.Println(msg)
    }
}

