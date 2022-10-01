package cluster

import "go-redis/interface/resp"

func execSelect(clusterDatabase *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return clusterDatabase.db.Exec(c, cmdArgs)
}
