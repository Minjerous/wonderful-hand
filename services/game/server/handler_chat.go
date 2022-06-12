package server

import (
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"wonderful-hand-game/network/protocol/packet"
)

type ChatHandler struct{}

func (*ChatHandler) Handle(pk packet.Packet, _ *Session) error {
	logging.Infof("%#v", pk)
	// todo broadcast to Room server
	return nil
}
