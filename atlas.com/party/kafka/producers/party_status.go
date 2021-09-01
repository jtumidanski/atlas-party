package producers

import "github.com/sirupsen/logrus"

type partyStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func PartyCreated(l logrus.FieldLogger) func(worldId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_STATUS")
	return func(worldId byte, partyId uint32, characterId uint32) {
		emitPartyStatus(producer, worldId, partyId, characterId, "CREATED")
	}
}

func PartyDisbanded(l logrus.FieldLogger) func(worldId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_STATUS")
	return func(worldId byte, partyId uint32, characterId uint32) {
		emitPartyStatus(producer, worldId, partyId, characterId, "DISBANDED")
	}
}

func emitPartyStatus(producer func(key []byte, event interface{}), worldId byte, partyId uint32, characterId uint32, theType string) {
	e := &partyStatusEvent{
		WorldId:     worldId,
		PartyId:     partyId,
		CharacterId: characterId,
		Type:        theType,
	}
	producer(CreateKey(int(partyId)), e)
}
