package packet

import "wonderful-hand-game/network/protocol"

// AddRoom 由客户端发出创建一个房间
type AddRoom struct {
	UID      uint64
	Token    string
	Password string
}

func (a *AddRoom) ID() uint8 {
	return IDAddRoom
}

func (a *AddRoom) Marshal(_ *protocol.Writer) {
}

func (a *AddRoom) Unmarshal(r *protocol.Reader) {
	r.VarUint64(&a.UID)
	r.String(&a.Token)
	r.String(&a.Password)
}
