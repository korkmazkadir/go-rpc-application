package main

import (
	"log"
	"net"
	"net/rpc"
	"time"

	"../../common"
	"../../peer"
	"../../registery"
)

type Node int

var nodeList []common.NodeInfo
var nodeInfo common.NodeInfo

var successor peer.Peer
var nodeRegistery registery.Registery

func (t *Node) Ping(counter int, reply *int) error {

	log.Println("Message arrived: ", counter)
	*reply = (counter + 1)

	go successor.Ping(counter + 1)

	return nil
}

func (t *Node) SendMessage(message string, reply *string) error {

	*reply = "ok"

	log.Println("Message arrived: ", message)

	return nil
}

func setSuccessor() {

	for idx, val := range nodeList {
		if val.PortNumber == nodeInfo.PortNumber {
			var successorNodeInfo = nodeList[(idx+1)%len(nodeList)]
			log.Println("successor is ", successorNodeInfo)
			successor.Init(successorNodeInfo)

			if idx == 0 {
				go func() {
					time.Sleep(5 * time.Second)
					//pingSuccessor(0)
					successor.Ping(0)
				}()
			}

			break
		}
	}

}

func (t *Node) SetNodeList(nodes []common.NodeInfo, reply *string) error {

	nodeList = nodes
	*reply = "ok"
	log.Println("Node list is assigned")
	setSuccessor()

	return nil
}

func main() {

	var node = new(Node)
	err := rpc.Register(node)

	if err != nil {
		log.Fatal("error registering nodeRegistery", err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		log.Fatal("error creating TCP address", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("error creating server port", err)
	}

	defer listener.Close()

	nodeInfo.IpAddress = "127.0.0.1"
	nodeInfo.PortNumber = listener.Addr().(*net.TCPAddr).Port
	nodeInfo.PublickKey = "ABC"

	nodeRegistery.Init("127.0.0.1", 8181)
	go nodeRegistery.Register(nodeInfo)

	log.Println("node is running...")
	for {
		rpc.Accept(listener)
	}

}
