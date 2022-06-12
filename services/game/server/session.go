package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"sync"
	"wonderful-hand-game/network"
	"wonderful-hand-game/network/protocol"
	"wonderful-hand-game/network/protocol/packet"
	"wonderful-hand-user/rpc/user"
)

type Session struct {
	s *Server

	rd *wsutil.Reader // rd 是 ws frame 的 reader
	wd *wsutil.Writer // wd 是 ws frame 的 writer

	data           *network.ClientData
	conn           gnet.Conn
	once, onceConn sync.Once
	handlers       map[uint8]PacketHandler // handler 用于处理 binary packet
}

func NewSession(srv *Server, conn gnet.Conn) *Session {
	s := &Session{
		s:        srv,
		conn:     conn,
		once:     sync.Once{},
		onceConn: sync.Once{},
		handlers: map[uint8]PacketHandler{},
	}
	s.registerHandlers()
	return s
}

var (
	errUnexpectHandshake = errors.New("unexpect handshake packet")
	errTokenInvalid      = errors.New("access token is invalid")
	errWriteFailed       = errors.New("write wd failed")
)

func (s *Session) Handshake() error {
	if s.rd == nil {
		s.rd = wsutil.NewServerSideReader(s.conn)
	}
	if s.wd == nil {
		s.wd = wsutil.NewWriter(s.conn, ws.StateServerSide, ws.OpBinary)
	}
	_, err := s.rd.NextFrame()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(s.rd)

	pid, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if pid != packet.IDOpenConnect1 {
		return errUnexpectHandshake
	}
	r := protocol.NewReader(reader)

	pk := &packet.OpenConnect1{}
	pk.Unmarshal(r)

	resp, err := s.s.rpcClis.UserRPC.UserTokenVerify(
		context.Background(),
		&user.UserTokenVerifyRequest{AccessToken: pk.Token},
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != 0 {
		return errTokenInvalid
	}

	s.data = &network.ClientData{
		UID:   pk.UID,
		Name:  pk.Name,
		Token: pk.Token,
	}

	return nil
}

func (s *Session) WriteBytes(b []byte) error {
	s.wd.Reset(s.conn, ws.StateServerSide, ws.OpBinary)
	n, err := s.wd.Write(b)
	if n != len(b) || err != nil {
		return errWriteFailed
	}
	return s.wd.Flush()
}

func (s *Session) WritePacket(pk packet.Packet) error {
	s.wd.Reset(s.conn, ws.StateServerSide, ws.OpBinary)
	buf := bufio.NewWriter(s.wd)
	pk.Marshal(protocol.NewWriter(buf))
	if err := buf.Flush(); err != nil {
		return err
	}
	return s.wd.Flush()
}

func (s *Session) registerHandlers() {
	s.handlers = map[uint8]PacketHandler{
		packet.IDChat:     &ChatHandler{},
		packet.IDMovePawn: &MovePawnHandler{},
	}
}

func (s *Session) GetData() *network.ClientData {
	return s.data
}

func (s *Session) ReadPacket() (packet.Packet, error) {
	if s.rd == nil {
		return nil, errors.New("this conn does not handshaked")
	}
	_, err := s.rd.NextFrame()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(s.rd)
	pid, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	pkFunc, ok := s.s.packetPool[pid]
	if !ok {
		return nil, errors.New("packet does not exist")
	}
	pk := pkFunc()
	pk.Unmarshal(protocol.NewReader(reader))
	return pk, nil
}

// HandlePacket 时 tcp 连接已经准备好了读取 (epoll)
func (s *Session) HandlePacket() error {
	pk, err := s.ReadPacket()
	if err != nil {
		return err
	}
	if err := s.handlePacket(pk); err != nil {
		logging.Debugf("err handling packet %#v, err %v", pk, err)
		return nil
	}
	return nil
}

func (s *Session) handlePacket(pk packet.Packet) error {
	handler, ok := s.handlers[pk.ID()]
	if !ok {
		logging.Errorf("cannot handle packet %#v, packet id %d, skip it", pk, pk.ID())
		return nil
	}
	if handler == nil {
		return nil // do nothing when handler is nil
	}
	if err := handler.Handle(pk, s); err != nil {
		return fmt.Errorf("%T: %w", pk, err)
	}
	return nil
}

func (s *Session) Close() error {
	s.once.Do(s.close)
	return nil
}

func (s *Session) close() {
	s.onceConn.Do(func() {
		s.s.CloseSession(s)
	})
	// todo 退出消息
}
