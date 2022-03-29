package character

import (
	"atlas-party/kafka"
	"atlas-party/party"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameStatus = "character_status_event"
	topicTokenStatus   = "TOPIC_CHARACTER_STATUS"
)

func StatusConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[statusEvent](consumerNameStatus, topicTokenStatus, groupId, handleStatus())
}

type statusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func handleStatus() kafka.HandlerFunc[statusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statusEvent) {
		if event.Type == "LOGIN" {
			party.MemberLogin(l, span)(event.CharacterId, event.WorldId, event.ChannelId)
		} else if event.Type == "LOGOUT" {
			party.MemberLogout(l, span)(event.CharacterId, event.WorldId, event.ChannelId)
		}
	}
}
