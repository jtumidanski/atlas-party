package consumers

import (
	"atlas-party/kafka/handler"
	"atlas-party/party"
	"github.com/sirupsen/logrus"
)

type partyExpelCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	PartyId     uint32 `json:"party_id"`
}

func EmptyPartyExpelCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &partyExpelCommand{}
	}
}

func HandlePartyExpelCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*partyExpelCommand); ok {
			party.Expel(l)(event.WorldId, event.ChannelId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
