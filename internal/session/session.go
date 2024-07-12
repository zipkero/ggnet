package session

import (
	"github.com/google/uuid"
	"net"
)

type Session struct {
	ID   string
	conn net.Conn
}

func NewSession(conn net.Conn) *Session {
	ss := &Session{conn: conn, ID: uuid.New().String()}
	return ss
}

func (s *Session) ReceiveMessages() {

}

func (s *Session) SendMessages() {

}
