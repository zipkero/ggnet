package host

import (
	"github.com/zipkero/ggnet/internal/handler"
	"github.com/zipkero/ggnet/internal/message"
	"github.com/zipkero/ggnet/internal/session"
	"log"
	"net"
	"sync"
)

type Host struct {
	endPoint *net.TCPAddr
	sessions map[string]*session.Session
	mu       sync.Mutex
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

func (h *Host) addSession(ss *session.Session) {
	h.sessions[ss.ID] = ss
}

func (h *Host) handleClient(conn net.Conn) {
	var sessionHandler handler.SessionHandler = h
	ss := session.NewSession(conn, sessionHandler)

	h.mu.Lock()
	h.addSession(ss)
	h.mu.Unlock()

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
}
