package service

import (
	"context"
	"wonderful-hand-game/rpc/game"
	"wonderful-hand-game/server"
)

type GameService struct {
	game.UnimplementedGameServiceServer
	s *server.Server
}

func New(s *server.Server) *GameService {
	return &GameService{s: s}
}

var (
	StatusBad int32 = 2
)

func (g *GameService) SendPacket(
	_ context.Context,
	req *game.GameSendPacketRequest,
) (resp *game.GameSendPacketResponse, _ error) {
	resp = new(game.GameSendPacketResponse)
	session, ok := g.s.GetSession(req.Uid)
	if !ok {
		resp.StatusCode = StatusBad
		resp.StatusMsg = "no session"
		return
	}
	err := session.WriteBytes(req.GetData())
	if err != nil {
		resp.StatusCode = StatusBad
		resp.StatusMsg = "write bytes failed"
		return
	}
	return
}
