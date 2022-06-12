package packet

import (
	"wonderful-hand-game/network/protocol"
)

type Packet interface {
	// ID 返回数据包的 ID
	ID() uint8
	// Marshal 将数据包序列化到 writer buffer 里
	Marshal(w *protocol.Writer)
	// Unmarshal 将 Reader buffer 内的数据包反序列化到 packet 内
	Unmarshal(r *protocol.Reader)
}
