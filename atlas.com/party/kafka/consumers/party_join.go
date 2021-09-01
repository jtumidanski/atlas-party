package consumers

import (
	"atlas-party/kafka/handler"
	"atlas-party/party"
	"github.com/sirupsen/logrus"
)

type partyJoinCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	PartyId     uint32 `json:"party_id"`
}

func EmptyPartyJoinCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &partyJoinCommand{}
	}
}

func HandlePartyJoinCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*partyJoinCommand); ok {
			party.Join(l)(event.WorldId, event.ChannelId, event.PartyId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}