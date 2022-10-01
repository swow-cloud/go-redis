package database

import (
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/lib/wildcard"
	"go-redis/resp/reply"
)

//DEL k1 k2 k3
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine2("del", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

//EXISTS
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exist := db.GetEntity(key)
		if exist {
			result++
		}

	}
	return reply.MakeIntReply(result)
}

//FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	db.addAof(utils.ToCmdLine2("flushdb", args...))
	return reply.MakeOkReply()
}

//TYPE k1
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none") //TCP:none\r\n
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
		//TODO 实现其他的数据类型
	}
	return &reply.UnknownErrReply{}
}

//KEYS

//RENAME k1 k2 k1:v k2:v
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeErrReply("no suck key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.addAof(utils.ToCmdLine2("rename", args...))
	return reply.MakeOkReply()
}

//RENAMENX k1 k2 如果k2存在就不覆盖，如果k2存在就删掉k1并且把k1的值保存到k2
func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	_, ok := db.GetEntity(dest)
	if ok {
		return reply.MakeIntReply(0)
	}
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeErrReply("no suck key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.addAof(utils.ToCmdLine2("renamenx", args...))
	return reply.MakeIntReply(1)
}

//KEYS *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.Foreach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key)) //[]byte(key) 转为字节切片
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}
func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("FlushDb", execFlushDB, -1)  //FLUSHDB a b c
	RegisterCommand("type", execType, 2)         //TYPE k1
	RegisterCommand("Rename", execRename, 3)     //RENAME k1 k2
	RegisterCommand("RenameNx", execRenameNx, 3) //RENAMENX k1 k2
	RegisterCommand("Keys", execKeys, 2)         //KEYS *

}
