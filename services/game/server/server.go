package server

import (
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"golang.org/x/sys/unix"
	"log"
	"sync"
	"sync/atomic"
	"wonderful-hand-game/internal/config"
	"wonderful-hand-game/network/protocol/packet"
	"wonderful-hand-user/rpc/user"
)

type Server struct {
	gnet.BuiltinEventEngine

	c         config.Config
	eng       gnet.Engine
	connected int64 // 连接数

	packetPool packet.Pool

	// p hold all sessions  uid -> Session
	p   map[uint64]*Session
	pmu sync.Mutex

	rpcClis *RPCClis
}

type RPCClis struct {
	UserRPC user.UserServiceClient
}

type ConnContext struct {
	ws         bool   // 是否已经升级协议
	handshaked bool   // 是否已经握手连接
	uid        uint64 // 持有 session 的 uid
}

func (s *Server) OnBoot(eng gnet.Engine) gnet.Action {
	s.eng = eng
	logging.Infof("server with multi-core=%t is listening on %s",
		s.c.Server.Multicore, s.c.Network.GameAddr)
	return gnet.None
}

func (s *Server) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(new(ConnContext))
	atomic.AddInt64(&s.connected, 1)
	return nil, gnet.None
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil && err != unix.ECONNRESET {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt64(&s.connected, -1)
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	if !c.Context().(*ConnContext).ws {
		// 升级协议
		_, err := ws.Upgrade(c)
		logging.Infof("conn[%v] upgrade websocket protocol", c.RemoteAddr().String())
		if err != nil {
			logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
			return gnet.Close
		}
		c.Context().(*ConnContext).ws = true
		return gnet.None
	} else if !c.Context().(*ConnContext).handshaked {
		// 连接握手
		se := NewSession(s, c)
		if se.Handshake() != nil {
			// 握手失败
			pk := &packet.OpenConnect2{StatusCode: packet.StatusCodeBad}
			_ = se.WritePacket(pk)
			return gnet.Close
		}

		// 握手成功，建立连接
		uid := se.GetData().UID
		s.pmu.Lock()
		s.p[uid] = se
		s.pmu.Unlock()

		c.Context().(*ConnContext).handshaked = true
		c.Context().(*ConnContext).uid = uid

		pk := &packet.OpenConnect2{StatusCode: packet.StatusCodeOK}
		if se.WritePacket(pk) != nil {
			return gnet.Close
		}
		return gnet.None
	}

	se, ok := s.p[c.Context().(*ConnContext).uid]
	if !ok {
		// 这里代表这个连接的 session 已经 Close 了
		return gnet.Close
	}

	// 处理数据包
	if se.HandlePacket() != nil {
		// 处理发生了错误
		return gnet.Close
	}

	return gnet.None
}

func (s *Server) Run() {
	log.Fatalln("server exits:",
		gnet.Run(
			s,
			s.c.Network.GameAddr,
			gnet.WithMulticore(s.c.Server.Multicore),
			gnet.WithReusePort(s.c.Server.ReusePort),
		),
	)
}

func (s *Server) CloseSession(se *Session) {
	if se.GetData() != nil {
		s.pmu.Lock()
		delete(s.p, se.GetData().UID)
		s.pmu.Unlock()
	}
}

func (s *Server) GetSession(uid uint64) (se *Session, ok bool) {
	se, ok = s.p[uid]
	return
}

func New(c *config.Config, clis *RPCClis) *Server {
	return &Server{
		c:          *c,
		p:          make(map[uint64]*Session),
		pmu:        sync.Mutex{},
		packetPool: packet.NewPool(),
		rpcClis:    clis,
	}
}

func Run(c *config.Config, clis *RPCClis) {
	New(c, clis).Run()
}
