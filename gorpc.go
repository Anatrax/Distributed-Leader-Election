package main

import (
    "fmt"
    "os"
	"log"
    "time"
	"strconv"
    "github.com/anatrax/better_grpc_base/node"
)

func parseArgs() (int, int) {
    //get command line arguments 
    num_nodes, err := strconv.Atoi(os.Args[1])
	if err!=nil {
		log.Fatal("error parsing command line arguments", err)
	}
    node_id, err := strconv.Atoi(os.Args[2])
	if err!=nil {
		log.Fatal("error parsing command line arguments", err)
	}
    return node_id, num_nodes
}

var this *node.Node
func main(){
    node_id, num_nodes := parseArgs()
    this = node.Init_node(node_id, num_nodes)
    //this.registerHandler(node.MessageType, pong)
    for i:= 0; i < num_nodes; i++ {
        this.Ping(i)
    }

    // Stall so I don't die before recieving all pings
    time.Sleep(5 * time.Second)
}

func pong(req_id int, res *int) error {
    fmt.Println("><><><><>")
//    this.Log("pong")
    return nil
}
