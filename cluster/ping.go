package cluster

import "go-redis/interface/resp"

func ping(clusterDatabase *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return clusterDatabase.db.Exec(c, cmdArgs)
}
