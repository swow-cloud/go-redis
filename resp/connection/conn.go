package connection

import (
	"go-redis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	selectedDB   int
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}

//返回Ip地址

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.mu.Unlock()
		c.waitingReply.Done()
	}()
	_, err := c.conn.Write(bytes)
	return err

}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.selectedDB = dbNum
}
