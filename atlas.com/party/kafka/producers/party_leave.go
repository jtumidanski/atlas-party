package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type partyLeaveCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func LeaveParty(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_PARTY_LEAVE")
	return func(worldId byte, channelId byte, characterId uint32) {
		e := &partyLeaveCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		}
		producer(CreateKey(int(characterId)), e)
	}
}