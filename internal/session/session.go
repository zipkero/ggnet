package session

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/zipkero/ggnet/internal/handler"
	"github.com/zipkero/ggnet/pkg/message"
	"io"
	"log"
	"net"
)

type Session struct {
	ID     string
	SendCh chan message.Message

	conn    net.Conn
	handler handler.SessionHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSession(conn net.Conn, handler handler.SessionHandler, ctx context.Context, cancel context.CancelFunc) *Session {
	ss := &Session{
		conn:    conn,
		ID:      uuid.New().String(),
		SendCh:  make(chan message.Message),
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
	return ss
}

func (s *Session) ReceiveMessages() {
	log.Println("ReceiveMessages started")
	defer log.Println("ReceiveMessages ended")

	readCh := make(chan message.Message)

	go func() {
		for {
			lengthBuffer := make([]byte, 4)
			_, err := io.ReadFull(s.conn, lengthBuffer)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println(err)
				continue
			}

			messageLength := binary.BigEndian.Uint32(lengthBuffer)
			messageBuffer := make([]byte, messageLength)

			_, err = io.ReadFull(s.conn, messageBuffer)
			if err != nil {
				log.Println(err)
				continue
			}

			messageType := binary.BigEndian.Uint16(messageBuffer[:2])
			messagePayload := messageBuffer[2:]

			readCh <- message.Message{
				Type:    messageType,
				Content: string(messagePayload),
			}
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-readCh:
			switch {
			case msg.Type > 10000:
				s.handler.HandleMessage(s.ID, msg)
			case msg.Type < 10000:
				s.HandleMessage(s.ID, msg)
			}
		}
	}
}

func (s *Session) SendMessages() {
	log.Println("SendMessages started")
	defer log.Println("SendMessages ended")

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.SendCh:
			var msgBytes []byte
			binary.BigEndian.PutUint16(msgBytes, msg.Type)
			msgBytes = append(msgBytes, []byte(msg.Content)...)

			length := uint32(len(msgBytes))
			lengthBuffer := new(bytes.Buffer)
			err := binary.Write(lengthBuffer, binary.BigEndian, length)
			if err != nil {
				log.Println(err)
				continue
			}

			sendMessage := append(lengthBuffer.Bytes(), msgBytes...)

			_, err = s.conn.Write(sendMessage)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (s *Session) HandleMessage(sessionId string, msg message.Message) {
	switch msg.Type {
	case 9999:
		log.Printf("received from: %s, type: %d, message: %s", sessionId, msg.Type, msg.Content)
	case 0:
		log.Printf("received from: %s, type: %d, message: %s", sessionId, msg.Type, msg.Content)
	default:
		log.Printf("received from: %s, type: %d, message: %s", sessionId, msg.Type, msg.Content)
	}
}

func (s *Session) Cancel() {
	s.cancel()
}
