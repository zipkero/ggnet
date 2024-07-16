package session

import (
	"bytes"
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
}

func NewSession(conn net.Conn, handler handler.SessionHandler) *Session {
	ss := &Session{
		conn:    conn,
		ID:      uuid.New().String(),
		sendCh:  make(chan message.Message),
		handler: handler,
	}
	return ss
}

func (s *Session) ReceiveMessages() {
	for {
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

		s.handler.HandleMessage(s.ID, message.Message{
			Type:    messageType,
			Content: string(messagePayload),
		})
	}
}

func (s *Session) SendMessages() {
	for {
		select {
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
