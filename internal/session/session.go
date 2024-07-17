package session

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/zipkero/ggnet/internal/handler"
	"github.com/zipkero/ggnet/internal/message"
	"io"
	"log"
	"net"
)

type Session struct {
	ID      string
	conn    net.Conn
	sendCh  chan message.Message
	handler handler.SessionHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSession(conn net.Conn, handler handler.SessionHandler, ctx context.Context, cancel context.CancelFunc) *Session {
	ss := &Session{
		conn:    conn,
		ID:      uuid.New().String(),
		sendCh:  make(chan message.Message),
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
	return ss
}

func (s *Session) ReceiveMessages() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			lengthBuffer := make([]byte, 4)
			_, err := io.ReadFull(s.conn, lengthBuffer)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println(err)
			}

			messageLength := binary.BigEndian.Uint32(lengthBuffer)
			messageBuffer := make([]byte, messageLength)

			_, err = io.ReadFull(s.conn, messageBuffer)
			if err != nil {
				log.Println(err)
			}

			messageType := binary.BigEndian.Uint16(messageBuffer[:2])
			messagePayload := messageBuffer[2:]

			switch {
			case messageType > 10000:
				s.handler.HandleMessage(s.ID, message.Message{
					Type:    messageType,
					Content: string(messagePayload),
				})
			case messageType < 10000:
				s.HandleMessage(s.ID, message.Message{
					Type:    messageType,
					Content: string(messagePayload),
				})
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

func (s *Session) SendMessages() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.sendCh:
			var msgBytes []byte
			binary.BigEndian.PutUint16(msgBytes, msg.Type)
			msgBytes = append(msgBytes, []byte(msg.Content)...)

			length := uint32(len(msgBytes))
			lengthBuffer := new(bytes.Buffer)
			err := binary.Write(lengthBuffer, binary.BigEndian, length)
			if err != nil {
				log.Println(err)
			}

			sendMessage := append(lengthBuffer.Bytes(), msgBytes...)

			_, err = s.conn.Write(sendMessage)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Session) Cancel() {
	s.cancel()
}
