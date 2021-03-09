package node

import (
	"fmt"
    "time"
	"log"
	"strconv"
	"net"
	"net/rpc"
	"net/http"
)

//we can add variables to this struct for node state
type Node struct {
    Id int // Node ID
    channels []*rpc.Client // Communication channels between other nodes
    Num_peers int // Number of other nodes
}

// Custom logging function 
func (node *Node) Log(str string){
	time := time.Now()
    // Print the time, node id, and provided message
	fmt.Printf("[%s | ID %03d]  %s\n", time.Format("15:04:05.0000"), node.Id, str)
}

// Node constructor
func Init_node(node_id int, total_nodes int) *Node {
	node := new(Node)
	node.Id = node_id
    node.Num_peers = total_nodes

	node.init_server()
    node.init_channels(total_nodes)
	node.Log("initialized server and connected to all channels")

	return node
}

//start listening for connections from other nodes
func (node *Node) init_server(){
    //register handlers
    //rpc register will export all functions with of type "Node" that have a captial first letter
    //e.g. Pong will be exported, ping will not
	rpc.Register(node)
	rpc.HandleHTTP()

    //start listing at (arbitrary) port 1200
	l, err := net.Listen("tcp", ":" + strconv.Itoa(1200 + node.Id))
	if err != nil {
		log.Fatal("could not initialize server", err)
	}
	go http.Serve(l, nil)
}

// Create connections to each other node in the network,
// and store their information in the clients array
func (node *Node) init_channels(total_nodes int){
    for i := 0; i < total_nodes; i++ {
        // Attempt to initiate connections with each node
		client, err := rpc.DialHTTP("tcp", "127.0.0.1:" + strconv.Itoa(1200 + i))
		if err != nil {
            // Retry indefinitely if connection fails
			for err!=nil{
				client, err = rpc.DialHTTP("tcp", "127.0.0.1:" + strconv.Itoa(1200 + i))
			}
		}
        str := fmt.Sprintf("opened channel to node %d", i)
        node.Log(str)
        node.channels = append(node.channels, client)
	}
	return
}

// Print a recieved message
func (node *Node) Ping(id int) int{
    var res int
    node.Log(fmt.Sprintf("pinging node %d", id))
    err := node.channels[id].Call("Node.Pong", node.Id, &res)
    if err != nil {
        fmt.Printf("unable to send ping\n", err)
    }
    return res
}

//receive a heartbeat from the leader, and update my last_com time
func (node *Node) Pong(id int, reply *int) error{
	//update last_com time to now
    node.Log("pong")
	return nil
}

