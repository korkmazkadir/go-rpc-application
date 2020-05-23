package peer

import (
	"fmt"
	"log"
	"net/rpc"

	"../common"
)

type Peer struct {
	nodeInfo common.NodeInfo
	client   *rpc.Client
}

func (t *Peer) Init(nodeInfo common.NodeInfo) {

	var address = fmt.Sprintf("%s:%d", nodeInfo.IpAddress, nodeInfo.PortNumber)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal("Peer connection error: ", err)
	}

	t.client = client
	t.nodeInfo = nodeInfo

}

func (t *Peer) Ping(counter int) {

	var result int
	err := t.client.Call("Node.Ping", counter, &result)
	if err != nil {
		log.Fatal("ping error: ", err)
	}

}
