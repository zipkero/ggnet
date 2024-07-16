package server

import "github.com/zipkero/ggnet/internal/host"

type Server struct {
	host *host.Host
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
