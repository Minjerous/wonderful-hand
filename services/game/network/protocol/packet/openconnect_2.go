package packet

import (
	"wonderful-hand-game/network/protocol"
)

const (
	StatusCodeOK uint8 = iota
	StatusCodeBad
)

// OpenConnect2 由服务端回复客户端的 OpenConnect1
type OpenConnect2 struct {
	StatusCode uint8
}

func (o *OpenConnect2) ID() uint8 {
	return IDOpenConnect2
}

func (o *OpenConnect2) Marshal(w *protocol.Writer) {
	w.Magic()
	w.Uint8(&o.StatusCode)
}

func (o *OpenConnect2) Unmarshal(_ *protocol.Reader) {
}
