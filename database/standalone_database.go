package database

import (
	"go-redis/aof"
	"go-redis/config"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

// Exec set k v
//get k
//select 2
func (database *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}
	dbIndex := client.GetDBIndex()
	db := database.dbSet[dbIndex]
	return db.Exec(client, args)
}

func (database *StandaloneDatabase) Close() {

}

func (database *StandaloneDatabase) AfterClientClose(client resp.Connection) {

}

func NewStandaloneDatabase() *StandaloneDatabase {
	database := &StandaloneDatabase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := MakeDB()
		db.index = i
		database.dbSet[i] = db
	}
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			////db=dbSet[0]
			////db=dbSet[1]
			////...
			////db=dbSet[15]
			//TODO 这里会造成一个问题，func闭包会引用外部变量db，每次循环虽然db的值会变， 但是内存地址不变，当最后一次循环时db值变为15，那么之前循环的func闭包里的db.index也会变为15
			//db.addAof = func(line CmdLine) {
			//	database.aofHandler.AddAof(db.index, line)
			//}

			sdb := db
			sdb.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(sdb.index, line)
			}
		}
	}
	return database
}

//select 2
func execSelect(connection resp.Connection, database *StandaloneDatabase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil { //select asnlknln
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) { //select 11111000022222
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	connection.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
