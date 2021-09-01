package producers

import "github.com/sirupsen/logrus"

type partyJoinCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	PartyId     uint32 `json:"party_id"`
}

func JoinParty(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_JOIN")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		e := &partyJoinCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			PartyId:     partyId,
			CharacterId: characterId,
		}
		producer(CreateKey(int(characterId)), e)
	}
}