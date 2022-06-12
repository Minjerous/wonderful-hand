package packet

import (
	protocol2 "wonderful-hand-game/network/protocol"
)

const (
	TextTypeRaw    uint8 = iota // debug only
	TextTypeChat                // 玩家对话
	TextTypeSystem              // 系统提示
)

type Chat struct {
	// TextType 定义消息的类型
	TextType uint8
	// SourceName 是消息的源，可以是用户名
	SourceName string
	// DestRoomID 房间 id
	DestRoomID string
	// Content 消息主体
	Content string
	// UID 用户的 id，只有当 TextType 为 TextTypeChat 才会定义
	UID uint64
}

func (c *Chat) ID() uint8 {
	return IDChat
}

func (c *Chat) Marshal(w *protocol2.Writer) {
	w.Uint8(&c.TextType)
	w.String(&c.SourceName)
	w.String(&c.DestRoomID)
	w.String(&c.Content)
	if c.TextType == TextTypeChat {
		w.VarUint64(&c.UID)
	}
}

func (c *Chat) Unmarshal(r *protocol2.Reader) {
	r.Uint8(&c.TextType)
	r.String(&c.SourceName)
	r.String(&c.DestRoomID)
	r.String(&c.Content)
	if c.TextType == TextTypeChat {
		r.VarUint64(&c.UID)
	}
}
