package consumers

import (
	"atlas-party/kafka/handler"
	"atlas-party/party"
	"github.com/sirupsen/logrus"
)

type partyCreateCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func EmptyPartyCreateCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &partyCreateCommand{}
	}
}

func HandlePartyCreateCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*partyCreateCommand); ok {
			party.Create(l)(event.CharacterId, event.WorldId, event.ChannelId)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
