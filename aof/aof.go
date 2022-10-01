package aof

import (
	"go-redis/config"
	"go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"os"
	"strconv"
)

type CmdLine = [][]byte

//var aofBufferSize = 1 << 16
const aofBufferSize = 1 << 16

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFileName string
	currentDb   int
}

//NewAofHandler

func NewAofHandler(database database.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFileName = config.Properties.AppendFilename
	handler.database = database
	//LoadAof
	handler.LoadAof()
	aofFile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	//channel
	handler.aofChan = make(chan *payload, aofBufferSize)
	go func() {
		handler.handleAof()
	}()
	return handler, nil
}

//Add payload(set k v) -> aofChan

func (handler *AofHandler) AddAof(dbIndex int, cmdLine CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}

}

//handleAof payload(set k v) <- aofChan(落盘)
func (handler *AofHandler) handleAof() {
	handler.currentDb = 0
	for p := range handler.aofChan { //读channel，但不确定channel是否关闭时使用for
		if p.dbIndex != handler.currentDb {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}

		}
		handler.currentDb = p.dbIndex
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}

	}
}

//LoadAof

func (handler *AofHandler) LoadAof() {

	file, err := os.Open(handler.aofFileName)
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()
	ch := parser.ParseStream(file)
	fackConn := &connection.Connection{}
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			logger.Error(p.Err)
			continue
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multi bulk")
			continue
		}
		rep := handler.database.Exec(fackConn, r.Args)
		if reply.IsErrReply(rep) {
			logger.Error("exec err", rep.ToBytes())
		}
	}
}
