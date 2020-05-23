package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"../../common"
)

type NodeRegistery int

var registeredNodes []common.NodeInfo

func (t *NodeRegistery) Register(node common.NodeInfo, reply *common.NodeInfo) error {

	registeredNodes = append(registeredNodes, node)
	*reply = node

	log.Println("A node registered: %+v", node)

	if len(registeredNodes) == 3 {
		go func() {
			log.Println("sending node list...")
			sendNodeList()
		}()
	}

	return nil
}

func (t *NodeRegistery) Unregister(node common.NodeInfo, reply *common.NodeInfo) error {

	for idx, val := range registeredNodes {
		if val.IpAddress == node.IpAddress {
			registeredNodes = append(registeredNodes[:idx], registeredNodes[idx+1:]...)
			*reply = node

			log.Println("A node unregistered: %+v", node)

			break
		}
	}

	return nil
}

func main() {

	var nodeRegistery = new(NodeRegistery)
	err := rpc.Register(nodeRegistery)

	if err != nil {
		log.Fatal("error registering nodeRegistery", err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8181")
	if err != nil {
		log.Fatal("error creating TCP address", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("error creating server port", err)
	}

	defer listener.Close()

	log.Println("node registery is running...")

	rpc.Accept(listener)

}

func sendList(ipAddress string, portNumber int) {

	var address = fmt.Sprintf("%s:%d", ipAddress, portNumber)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal("dialing error:", err)
	}

	var result string
	err = client.Call("Node.SetNodeList", registeredNodes, &result)
	if err != nil {
		log.Fatal("Node could not registered", err)
	}

}

func sendNodeList() {

	for _, val := range registeredNodes {
		sendList(val.IpAddress, val.PortNumber)
	}

}
