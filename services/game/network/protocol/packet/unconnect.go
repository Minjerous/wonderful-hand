package packet

import (
	"wonderful-hand-game/network/protocol"
)

// UnConnectRoomPing sent by client to query room info and response with UnConnectRoomPong
type UnConnectRoomPing struct {
	Magic []byte
	// DestRoomID room id
	DestRoomID string
}

func (u *UnConnectRoomPing) ID() uint8 {
	return IDUnConnectRoomPing
}

func (u *UnConnectRoomPing) Marshal(_ *protocol.Writer) {
}

func (u *UnConnectRoomPing) Unmarshal(r *protocol.Reader) {
	r.Magic()
	r.String(&u.DestRoomID)
}

const (
	RoomStatusRunning uint8 = iota
	RoomStatusWait
	RoomStatusEnd
)

// UnConnectRoomPong sent by server to response client's UnConnectRoomPing
type UnConnectRoomPong struct {
	RoomName     string
	RoomSubtitle string
	RoomStatus   uint8 // todo add more
}

func (u *UnConnectRoomPong) ID() uint8 {
	return IDUnConnectRoomPong
}

func (u *UnConnectRoomPong) Marshal(w *protocol.Writer) {
	w.Magic()
	w.String(&u.RoomName)
	w.String(&u.RoomSubtitle)
	w.Uint8(&u.RoomStatus)
}

func (u *UnConnectRoomPong) Unmarshal(_ *protocol.Reader) {
}
