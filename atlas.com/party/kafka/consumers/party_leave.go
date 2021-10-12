package consumers

import (
	"atlas-party/kafka/handler"
	"atlas-party/party"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type partyLeaveCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func EmptyPartyLeaveCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &partyLeaveCommand{}
	}
}

func HandlePartyLeaveCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*partyLeaveCommand); ok {
			party.Leave(l, span)(event.WorldId, event.ChannelId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
