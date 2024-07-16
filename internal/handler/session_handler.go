package handler

import (
	"github.com/zipkero/ggnet/internal/message"
)

type SessionHandler interface {
	HandleMessage(sessionId string, msg message.Message)
}
