package acceptor

import (
	"github.com/zipkero/ggnet/internal/session"
	"net"
	"sync"
)

type Acceptor struct {
	endPoint *net.TCPAddr
	sessions map[string]*session.Session
}

func NewAcceptor(endPoint string) (*Acceptor, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", endPoint)
	if err != nil {
		return nil, err
	}
	return &Acceptor{
		endPoint: tcpAddr,
		sessions: make(map[string]*session.Session),
	}, nil
}

func (a *Acceptor) Listen() {
	l, err := net.ListenTCP("tcp", a.endPoint)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		ss := session.NewSession(conn)
		a.sessions[ss.ID] = ss
		if err != nil {
			panic(err)
		}

		go a.handleClient(ss)
	}
}

func (a *Acceptor) handleClient(ss *session.Session) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		ss.SendMessages()
	}()

	go func() {
		defer wg.Done()

		ss.ReceiveMessages()
	}()

	wg.Wait()
}
