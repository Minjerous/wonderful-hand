package server

import (
	"wonderful-hand-game/network/protocol/packet"
)

type PacketHandler interface {
	Handle(p packet.Packet, s *Session) error
}
