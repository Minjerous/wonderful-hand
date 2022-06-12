package handshake

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"time"
	"wonderful-hand-game/network/protocol"
	"wonderful-hand-game/network/protocol/packet"
)

// handshake packet
func packet1() []byte {
	c1 := packet.OpenConnect1{
		UID:   114514,
		Name:  "iGxnon",
		Token: "fuck.dwai.dawiodj",
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	w := protocol.NewWriter(buf)
	w.Magic()
	w.VarUint64(&c1.UID)
	w.String(&c1.Name)
	w.String(&c1.Token)
	return buf.Bytes()
}

// chat packet
func packet2() []byte {
	c2 := packet.Chat{
		TextType:   packet.TextTypeChat,
		SourceName: "iGxnon",
		DestRoomID: "adwh.awd-0",
		Content:    "Hello",
		UID:        114514,
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	w := protocol.NewWriter(buf)
	c2.Marshal(w)
	return buf.Bytes()
}

func Handshake(addr string) {
	conn, _, _, err := ws.Dial(context.Background(), "ws://127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	err = wsutil.WriteClientBinary(conn, append([]byte{packet.IDOpenConnect1}, packet1()...))
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
	data, code, err := wsutil.ReadServerData(conn)
	fmt.Println(data, code, err)

	err = wsutil.WriteClientBinary(conn, append([]byte{packet.IDChat}, packet2()...))
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)

}
