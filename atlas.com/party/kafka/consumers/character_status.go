package consumers

import (
	"atlas-party/kafka/handler"
	"atlas-party/party"
	"github.com/sirupsen/logrus"
)

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func EmptyCharacterStatusEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterStatusEvent{}
	}
}

func HandleCharacterStatusEvent() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			if event.Type == "LOGIN" {
				party.MemberLogin(l)(event.CharacterId, event.WorldId, event.ChannelId)
			} else if event.Type == "LOGOUT" {
				party.MemberLogout(l)(event.CharacterId, event.WorldId, event.ChannelId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
