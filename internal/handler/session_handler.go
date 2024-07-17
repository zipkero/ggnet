package handler

import (
	"github.com/zipkero/ggnet/pkg/message"
)

type SessionHandler interface {
	HandleMessage(sessionId string, msg message.Message)
}
