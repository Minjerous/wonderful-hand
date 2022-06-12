package packet

import (
	"wonderful-hand-game/network/protocol"
)

// OpenConnect1 由客户端发送请求连接
type OpenConnect1 struct {
	UID   uint64
	Name  string
	Token string // 登录后给的 jwt
}

func (o *OpenConnect1) ID() uint8 {
	return IDOpenConnect1
}

func (o *OpenConnect1) Marshal(_ *protocol.Writer) {
}

func (o *OpenConnect1) Unmarshal(r *protocol.Reader) {
	r.Magic()
	r.VarUint64(&o.UID)
	r.String(&o.Name)
	r.String(&o.Token)
}
