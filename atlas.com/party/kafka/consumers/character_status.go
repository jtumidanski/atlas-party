package consumers

import (
	"atlas-party/kafka/handler"
	"atlas-party/party"
	"github.com/opentracing/opentracing-go"
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
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			if event.Type == "LOGIN" {
				party.MemberLogin(l, span)(event.CharacterId, event.WorldId, event.ChannelId)
			} else if event.Type == "LOGOUT" {
				party.MemberLogout(l, span)(event.CharacterId, event.WorldId, event.ChannelId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
