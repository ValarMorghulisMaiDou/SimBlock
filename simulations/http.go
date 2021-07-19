package simulations

import (
	"SimBlock/protocol"
	"SimBlock/simulations/adapters"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/gin-gonic/gin"
)

var allNetworks = make(map[string]*Network)

var services = map[string]adapters.LifecycleConstructor{
	"pingpong": func(ctx *adapters.ServiceContext, stack *node.Node) (node.Lifecycle, error) {
		stack.RegisterProtocols([]p2p.Protocol{protocol.NewPingPong()})
		return nil, nil
	},
}

var simAdapter = adapters.NewSimAdapter(services)

func newNetwork(c *gin.Context) {
	name := c.Query("name")
	if len(name) == 0 {
		c.JSON(http.StatusInternalServerError, "name can not be empty")
		return
	}
	if allNetworks[name] != nil {
		c.JSON(http.StatusInternalServerError, "network already exists")
		return
	}
	network := NewNetwork(simAdapter, &NetworkConfig{
		DefaultService: "pingpong",
	})
	allNetworks[name] = network
	c.JSON(http.StatusOK, name)
}

func delNetwork(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusInternalServerError, "need network name")
		return
	}
	network := allNetworks[name]
	if network == nil {
		c.JSON(http.StatusInternalServerError, "network doesn't exist")
		return
	}
	network.Shutdown()
	delete(allNetworks, name)
	c.JSON(http.StatusOK, name)
}

type networkInfo struct {
	Name       string `json:"name"`
	NodeNumber int    `json:"nodeNumber"`
	NodeType   string `json:"nodeType"`
}

// 获取当前所有运行的仿真网络
func networks(c *gin.Context) {
	result := []networkInfo{}
	for k, network := range allNetworks {
		info := networkInfo{
			Name:       k,
			NodeNumber: len(network.GetNodes()),
			NodeType:   "memory",
		}
		result = append(result, info)
	}
	c.JSON(http.StatusOK, result)
}

func startNetwork(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusInternalServerError, "need network name")
		return
	}
	network := allNetworks[name]
	if network == nil {
		c.JSON(http.StatusInternalServerError, "network doesn't exist")
		return
	}
	if err := network.StartAll(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, name)
}

func stopNetwork(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusInternalServerError, "need network name")
		return
	}
	network := allNetworks[name]
	if network == nil {
		c.JSON(http.StatusInternalServerError, "network doesn't exist")
		return
	}
	if err := network.StopAll(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, name)
}

func createNode(c *gin.Context) {
	config := &adapters.NodeConfig{}

	// 节点的配置以json的格式保存在请求体中
	if err := c.ShouldBind(config); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	network := parseNetwork(c)
	if network == nil {
		return
	}
	privk, err := crypto.GenerateKey()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	config.PrivateKey = privk
	config.ID = enode.PubkeyToIDV4(&privk.PublicKey)
	if len(config.Name) == 0 {
		config.Name = fmt.Sprintf("node%02d", len(network.GetNodes()))
	}
	node, err := network.NewNodeWithConfig(config)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, node.NodeInfoWithConns())
}
func getNodes(c *gin.Context) {
	network := parseNetwork(c)
	if network == nil {
		return
	}
	nodes := network.GetNodes()

	infos := make([]*SimNodeInfo, len(nodes))
	for i, node := range nodes {
		infos[i] = node.NodeInfoWithConns()
	}

	c.JSON(http.StatusOK, infos)
}
func getNode(c *gin.Context) {
	node, _ := parseNode(c)
	if node == nil {
		return
	}
	c.JSON(http.StatusOK, node.NodeInfoWithConns())
}

func delNode(c *gin.Context) {
	node, network := parseNode(c)
	if node == nil {
		return
	}
	if err := network.DeleteNode(node.ID()); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func startNode(c *gin.Context) {
	node, network := parseNode(c)
	if node == nil {
		return
	}

	if err := network.Start(node.ID()); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, node.NodeInfo())
}
func stopNode(c *gin.Context) {
	node, network := parseNode(c)
	if node == nil {
		return
	}

	if err := network.Stop(node.ID()); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, node.NodeInfo())
}
func connectNode(c *gin.Context) {
	node, peer, network := parsePeer(c)
	if peer == nil {
		return
	}
	if err := network.Connect(node.ID(), peer.ID()); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, node.NodeInfo())
}
func disconnectNode(c *gin.Context) {
	node, peer, network := parsePeer(c)
	if peer == nil {
		return
	}
	if err := network.Disconnect(node.ID(), peer.ID()); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, node.NodeInfo())
}
func nodeRPC(c *gin.Context) {
	network := parseNetwork(c)
	if network == nil {
		return
	}
}

func parseNetwork(c *gin.Context) *Network {
	networkName := c.Param("network")
	if len(networkName) == 0 {
		c.String(http.StatusBadRequest, "need network name")
		return nil
	}
	network := allNetworks[c.Param("network")]
	if network == nil {
		c.String(http.StatusBadRequest, "network doesn't exist")
		return nil
	}
	return network
}

func parseNode(c *gin.Context) (*Node, *Network) {
	network := parseNetwork(c)
	if network == nil {
		return nil, nil
	}
	strid := c.Param("nodeid")
	if strid == "" {
		c.JSON(http.StatusBadRequest, "need nodeid")
		return nil, nil
	}
	id, err := enode.ParseID(strid)
	if err != nil {
		c.JSON(http.StatusBadRequest, "wrong node id format")
		return nil, nil
	}
	node := network.GetNode(id)
	if node == nil {
		c.JSON(http.StatusBadRequest, "node doesn't exist")
		return nil, nil
	}
	return node, network
}

func parsePeer(c *gin.Context) (*Node, *Node, *Network) {
	node, network := parseNode(c)
	if node == nil {
		return nil, nil, nil
	}
	peerstrid := c.Param("peerid")
	peerid, err := enode.ParseID(peerstrid)
	if err != nil {
		c.JSON(http.StatusBadRequest, "wrong peer id format")
		return nil, nil, nil
	}
	peer := network.GetNode(peerid)
	if peer == nil {
		c.JSON(http.StatusBadRequest, "peer doesn't exist")
		return nil, nil, nil
	}
	return node, peer, network
}
