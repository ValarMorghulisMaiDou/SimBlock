package simulations

import "github.com/gin-gonic/gin"
import "SimBlock/config"

func RegisterSimnet(g *gin.Engine) {
	simnetGroup := g.Group(config.SIMNET_PREFIX)
	simnetGroup.GET("/newnetwork", newNetwork)
	simnetGroup.GET("/delnetwork", delNetwork)
	simnetGroup.GET("/networks", networks)

	simnetGroup.POST("/start", startNetwork)
	simnetGroup.POST("/stop", stopNetwork)

	nodeGroup := simnetGroup.Group("/nodes/:network")
	nodeGroup.POST("", createNode)
	nodeGroup.GET("", getNodes)
	nodeGroup.GET("/:nodeid", getNode)
	nodeGroup.DELETE("/:nodeid", delNode)
	nodeGroup.POST("/:nodeid/start", startNode)
	nodeGroup.POST("/:nodeid/stop", stopNode)
	nodeGroup.POST("/:nodeid/conn/:peerid", connectNode)
	nodeGroup.DELETE("/:nodeid/conn/:peerid", disconnectNode)
	nodeGroup.GET("/:nodeid/rpc", nodeRPC)
}
