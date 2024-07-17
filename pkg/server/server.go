package server

import (
	"github.com/zipkero/ggnet/internal/host"
	"github.com/zipkero/ggnet/internal/session"
	"github.com/zipkero/ggnet/pkg/message"
)

type Server struct {
	host Host
}

type Host interface {
	Listen() error
	UniCast(sessionId string, msg message.Message) error
	BroadCast(msg message.Message)
	KickSession(sessionId string) error
	GetSessions() []*session.Session
	FindSession(sessionId string) (*session.Session, error)
}

func NewServer(endPoint string) (*Server, error) {
	h, err := host.NewHost(endPoint)
	if err != nil {
		return nil, err
	}
	return &Server{host: h}, nil
}

func (s *Server) Start() error {
	return s.host.Listen()
}

func (s *Server) Kick(sessionId string) error {
	return s.host.KickSession(sessionId)
}

func (s *Server) UniCast(sessionId string, msg message.Message) error {
	return s.host.UniCast(sessionId, msg)
}

func (s *Server) BroadCast(msg message.Message) {
	s.host.BroadCast(msg)
}

func (s *Server) GetSessions() []*SessionInfo {
	sessions := s.host.GetSessions()
	sessionInfos := make([]*SessionInfo, 0, len(sessions))
	for _, ss := range sessions {
		sessionInfos = append(sessionInfos, &SessionInfo{
			ID: ss.ID,
		})
	}
	return sessionInfos
}

func (s *Server) FindSession(sessionId string) (*SessionInfo, error) {
	ss, err := s.host.FindSession(sessionId)
	if err != nil {
		return nil, err
	}
	return &SessionInfo{
		ID: ss.ID,
	}, nil
}
