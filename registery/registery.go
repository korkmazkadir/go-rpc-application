package registery

import (
	"fmt"
	"log"
	"net/rpc"

	"../common"
)

type Registery struct {
	ipAddress  string
	portNumber int
	client     *rpc.Client
}

func (t *Registery) Init(ipAddress string, portNumber int) {

	var address = fmt.Sprintf("%s:%d", ipAddress, portNumber)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal("Registery dialing error: ", err)
	}

	t.client = client
	t.ipAddress = ipAddress
	t.portNumber = portNumber

}

func (t *Registery) Register(nodeInfo common.NodeInfo) {

	var reply common.NodeInfo
	var err = t.client.Call("NodeRegistery.Register", nodeInfo, &reply)
	if err != nil {
		log.Fatal("Node could not be registered", err)
	}

	log.Println("Node is registered", reply)

}
