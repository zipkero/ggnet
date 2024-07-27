package host

import (
	"context"
	"github.com/zipkero/ggnet/internal/handler"
	"github.com/zipkero/ggnet/internal/session"
	"github.com/zipkero/ggnet/pkg/message"
	"log"
	"net"
	"sync"
)

type Host struct {
	endPoint *net.TCPAddr
	sessions map[string]*session.Session
	mu       sync.RWMutex
}

func NewHost(endPoint string) (*Host, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", endPoint)
	if err != nil {
		return nil, err
	}
	return &Host{
		endPoint: tcpAddr,
		sessions: make(map[string]*session.Session),
	}, nil
}

func (h *Host) Listen() error {
	l, err := net.ListenTCP("tcp", h.endPoint)
	if err != nil {
		return err
	}
	defer l.Close()

	log.Printf("Listening on %s", h.endPoint)

	for {
		var conn net.Conn
		conn, err = l.Accept()

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		go h.handleClient(conn)
	}
}

func (h *Host) HandleMessage(sessionId string, msg message.Message) {
	switch msg.Type {
	case 10000:
		log.Printf("received from: %s, type: %d, message: %s", sessionId, msg.Type, msg.Content)
	default:
		log.Printf("received from: %s, type: %d, message: %s", sessionId, msg.Type, msg.Content)
	}
}

func (h *Host) KickSession(sessionId string) error {
	ss, err := h.getSession(sessionId)
	if err != nil {
		return err
	}

	ss.Cancel()
	h.removeSession(sessionId)
	return nil
}

func (h *Host) UniCast(sessionId string, msg message.Message) error {
	ss, err := h.getSession(sessionId)
	if err != nil {
		return err
	}
	ss.SendCh <- msg
	return nil
}

func (h *Host) BroadCast(msg message.Message) {
	h.mu.RLock()
	sessions := make(map[string]*session.Session, len(h.sessions))
	for k, v := range h.sessions {
		sessions[k] = v
	}
	h.mu.RUnlock()

	// TODO: need to worker pool
	for _, ss := range sessions {
		go func(ss *session.Session) {
			h.mu.Lock()
			ss.SendCh <- msg
			h.mu.Unlock()
		}(ss)
	}
}

func (h *Host) GetSessions() []*session.Session {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessions := make([]*session.Session, len(h.sessions))
	for _, v := range h.sessions {
		sessions = append(sessions, v)
	}
	return sessions
}

func (h *Host) FindSession(sessionId string) (*session.Session, error) {
	return h.getSession(sessionId)
}

func (h *Host) handleClient(conn net.Conn) {
	var sessionHandler handler.SessionHandler = h

	ctx, cancel := context.WithCancel(context.Background())
	ss := session.NewSession(conn, sessionHandler, ctx, cancel)

	h.addSession(ss)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ss.ReceiveMessages()
	}()

	go func() {
		defer wg.Done()
		go ss.SendMessages()
	}()

	wg.Wait()

	err := conn.Close()
	if err != nil {
		log.Println(err)
	}

	h.removeSession(ss.ID)
}

func (h *Host) addSession(ss *session.Session) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sessions[ss.ID] = ss
}

func (h *Host) removeSession(sessionId string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.sessions, sessionId)
}

func (h *Host) getSession(sessionId string) (*session.Session, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ss, ok := h.sessions[sessionId]; ok {
		return ss, nil
	}
	return nil, session.ErrSessionNotFound
}
