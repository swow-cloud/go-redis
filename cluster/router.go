package cluster

import (
	"go-redis/interface/resp"
)

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc //exists k1
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["ping"] = ping
	routerMap["rename"] = Rename
	routerMap["renamenx"] = Rename
	routerMap["flushdb"] = flushDB
	routerMap["del"] = Del
	routerMap["select"] = execSelect
	return routerMap
}

//GET Key //Set K1 v1
func defaultFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.PickNode(key)
	return cluster.reply(peer, c, cmdArgs)
}
