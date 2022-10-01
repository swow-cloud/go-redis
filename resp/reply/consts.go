package reply

type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p PongReply) ToBytes() []byte {
	//TODO implement me
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

// OkReply is +OK
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

// ToBytes marshal redis.Reply
func (r *OkReply) ToBytes() []byte {
	return okBytes
}

var theOkReply = new(OkReply)

// MakeOkReply returns a ok reply
func MakeOkReply() *OkReply {
	return theOkReply
}

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

var emptyMultiBulkBytes = []byte("*0\r\n")

// EmptyMultiBulkReply is a empty list
type EmptyMultiBulkReply struct{}

// ToBytes marshal redis.Reply
func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type NoReply struct {
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}

var noBytes = []byte("")

func (n NoReply) ToBytes() []byte {
	return noBytes
}
