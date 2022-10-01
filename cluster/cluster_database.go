package cluster

import (
	"context"
	"fmt"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/config"
	database2 "go-redis/database"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/consistenthash"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"runtime/debug"
	"strings"
)

type ClusterDatabase struct {
	self           string
	nodes          []string
	peerPicker     *consistenthash.NodeMap
	peerConnection map[string]*pool.ObjectPool
	db             database.Database
}

func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, cluster.self)
	cluster.peerPicker.AddNodes(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	fmt.Println(cluster.peerConnection)
	cluster.nodes = nodes
	return cluster

}

type CmdFunc func(clusterDatabase *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &reply.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("not supported cmd")
	}
	result = cmdFunc(cluster, client, args)
	return
}

func (cluster *ClusterDatabase) Close() {
	cluster.db.Close()
}

func (cluster *ClusterDatabase) AfterClientClose(client resp.Connection) {
	cluster.db.AfterClientClose(client)
}
