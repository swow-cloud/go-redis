package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

type EchoDataBase struct {
}

func NewEchoDataBase() *EchoDataBase {
	return &EchoDataBase{}
}

func (e EchoDataBase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e EchoDataBase) Close() {
}

func (e EchoDataBase) AfterClientClose(client resp.Connection) {
}
