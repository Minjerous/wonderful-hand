package packet

var registeredPackets = map[uint8]func() Packet{}

// Pool todo 用 sync.Pool
type Pool map[uint8]func() Packet

func Register(id uint8, pk func() Packet) {
	registeredPackets[id] = pk
}

func NewPool() Pool {
	p := Pool{}
	for id, pk := range registeredPackets {
		p[id] = pk
	}
	return p
}

func init() {
	buildinPackets := map[uint8]func() Packet{
		IDChat:              func() Packet { return &Chat{} },
		IDOpenConnect1:      func() Packet { return &OpenConnect1{} },
		IDOpenConnect2:      func() Packet { return &OpenConnect2{} },
		IDUnConnectRoomPing: func() Packet { return &UnConnectRoomPing{} },
		IDUnConnectRoomPong: func() Packet { return &UnConnectRoomPong{} },
		IDMovePawn:          func() Packet { return &MovePawn{} },
	}

	for id, pk := range buildinPackets {
		Register(id, pk)
	}
}
