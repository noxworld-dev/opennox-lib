package noxnet

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/netip"
	"reflect"
	"sync"
)

type PlayerID uint32

type Player interface {
	PlayerID() PlayerID
	PlayerName() string
	Disconnect()
}

type Mapper interface {
	NewPlayer(addr netip.AddrPort, cli Player, p Player) (StreamID, error)
	GetPlayer(addr netip.AddrPort, sid StreamID) Player
	GetPlayerPeer(addr netip.AddrPort, cli Player, sid StreamID) Player
	DelPlayer(addr netip.AddrPort, cli Player, sid StreamID) bool
}

func NewMapper() Mapper {
	return &defaultMapper{}
}

type defaultMapper struct {
	mu    sync.RWMutex
	bySID [MaxStreams - 2]Player
}

func (m *defaultMapper) NewPlayer(_ netip.AddrPort, _ Player, p Player) (StreamID, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, p2 := range m.bySID {
		if p2 == nil {
			m.bySID[i] = p
			return StreamID(i + 1), nil
		}
	}
	return 0, ErrFull
}

func (m *defaultMapper) GetPlayer(_ netip.AddrPort, sid StreamID) Player {
	if sid == ServerStreamID || sid == MaxStreamID {
		return nil
	}
	ind := int(sid) - 1
	m.mu.RLock()
	p := m.bySID[ind]
	m.mu.RUnlock()
	return p
}

func (m *defaultMapper) GetPlayerPeer(addr netip.AddrPort, _ Player, sid StreamID) Player {
	return m.GetPlayer(addr, sid)
}

func (m *defaultMapper) DelPlayer(_ netip.AddrPort, _ Player, sid StreamID) bool {
	if sid == ServerStreamID || sid == MaxStreamID {
		return false
	}
	ind := int(sid) - 1
	m.mu.Lock()
	p := m.bySID[ind]
	m.bySID[ind] = nil
	m.mu.Unlock()
	return p != nil
}

type Engine interface {
	ServerInfo(addr netip.AddrPort) *MsgServerInfo
	PreJoin(addr netip.AddrPort, req *MsgServerTryJoin) error
	CheckPass(addr netip.AddrPort, pass string) error
	Connect(addr netip.AddrPort) (Player, error)
}

func NewServer(log *slog.Logger, conn PacketConn, e Engine, opts *ServerOptions) *Server {
	return NewServerWithPort(log, NewPort(log, conn, true), e, opts)
}

type ServerOptions struct {
	PlayerMap Mapper
	NoXor     bool
}

func NewServerWithPort(log *slog.Logger, port *Port, e Engine, opts *ServerOptions) *Server {
	if opts == nil {
		opts = &ServerOptions{}
	}
	if opts.PlayerMap == nil {
		opts.PlayerMap = NewMapper()
	}
	s := &Server{
		log:  log,
		e:    e,
		Port: port,
	}
	s.players.mapper = opts.PlayerMap
	s.players.noXor = opts.NoXor
	s.Port.OnMessage(s.handleMsg)
	s.Port.Start()
	return s
}

type Server struct {
	log  *slog.Logger
	Port *Port
	e    Engine

	players struct {
		mapper Mapper
		noXor  bool
	}
}

func (s *Server) LocalAddr() netip.AddrPort {
	return s.Port.LocalAddr()
}

func (s *Server) Close() {
	// TODO: disconnect all clients
	s.Reset()
	s.Port.Close()
}

func (s *Server) Reset() {
	s.Port.Reset()
}

func (s *Server) handleMsg(conn *Conn, sid StreamID, m Message) bool {
	srv := conn.WithID(ServerStreamID)
	switch sid {
	default:
		p := s.players.mapper.GetPlayer(conn.RemoteAddr(), sid)
		if p == nil {
			s.log.Warn("unhandled player message", "sid", sid, "type", reflect.TypeOf(m).String(), "msg", m)
			return false
		}
		return s.handlePlayerMsg(srv, p, m)
	case ServerStreamID:
		return s.handleGlobalMsg(srv, m)
	case MaxStreamID:
		return s.handleConnectMsg(srv, m)
	}
}

func (s *Server) handleGlobalMsg(conn *Stream, m Message) bool {
	switch m := m.(type) {
	default:
		s.log.Warn("unhandled global message", "type", reflect.TypeOf(m).String(), "msg", m)
		return false
	case *MsgDiscover:
		info := s.e.ServerInfo(conn.Addr())
		if info == nil {
			return true // ignore
		}
		info.Token = m.Token
		_ = conn.SendUnreliableMsg(info)
		return true
	case *MsgServerTryJoin:
		err := s.e.PreJoin(conn.Addr(), m)
		if e, ok := ErrorToMsg(err); ok && e != nil {
			_ = conn.SendUnreliableMsg(e)
			return true
		} else if err != nil {
			s.log.Error("cannot check join", "err", err)
			_ = conn.SendUnreliableMsg(&MsgJoinFailed{})
			return true
		}
		_ = conn.SendUnreliableMsg(&MsgJoinOK{})
		return true
	case *MsgServerPass:
		err := s.e.CheckPass(conn.Addr(), m.Pass)
		if e, ok := ErrorToMsg(err); ok && e != nil {
			_ = conn.SendUnreliableMsg(e)
			return true
		} else if err != nil {
			s.log.Error("cannot check pass", "err", err)
			_ = conn.SendUnreliableMsg(&MsgServerError{Err: ErrWrongPassword})
			return true
		}
		_ = conn.SendUnreliableMsg(&MsgJoinOK{})
		return true
	}
}

func (s *Server) handleConnectMsg(conn *Stream, m Message) bool {
	ctx := context.Background()
	switch m := m.(type) {
	case *MsgConnect:
		addr := conn.Addr()
		log := s.log.With("addr", addr)
		p, err := s.e.Connect(addr)
		if err != nil {
			log.Error("cannot connect player", "err", err)
			return true
		}
		sid, err := s.players.mapper.NewPlayer(addr, p, p)
		if err != nil {
			// TODO: send error to the client
			log.Error("cannot create player", "err", err)
			return true
		}
		var xor byte
		if !s.players.noXor {
			xor = byte(rand.UintN(0xff))
		}
		_, err = conn.QueueReliableMsg(ctx, []Message{&MsgAccept{
			ID: 0, // TODO
		}, &MsgServerAccept{
			ID: uint32(sid), XorKey: xor,
		}}, nil, func() {
			s.players.mapper.DelPlayer(addr, p, sid)
			p.Disconnect()
		})
		if err != nil {
			log.Error("cannot connect player", "err", err)
			return true
		}
		return true
	default:
		s.log.Warn("unhandled connect message", "type", reflect.TypeOf(m).String(), "msg", m)
		return false
	}
}

func (s *Server) handlePlayerMsg(conn *Stream, p Player, m Message) bool {
	switch m := m.(type) {
	default:
		s.log.Warn("unhandled player message", "player", p.PlayerID(), "type", reflect.TypeOf(m).String(), "msg", m)
		return false
	}
}
