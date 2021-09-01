package producers

import "github.com/sirupsen/logrus"

type partyMemberStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func PartyMemberLogin(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "LOGIN")
	}
}

func PartyMemberLogout(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "LOGOUT")
	}
}

func PartyMemberJoin(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "JOINED")
	}
}

func PartyMemberLeave(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "LEFT")
	}
}

func PartyMemberExpelled(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "EXPELLED")
	}
}

func PartyMemberDisbanded(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "DISBANDED")
	}
}

func PartyMemberPromoted(l logrus.FieldLogger) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "PROMOTED")
	}
}

func emitPartyMemberStatus(producer func(key []byte, event interface{}), worldId byte, channelId byte, partyId uint32, characterId uint32, theType string) {
	e := &partyMemberStatusEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		PartyId:     partyId,
		CharacterId: characterId,
		Type:        theType,
	}
	producer(CreateKey(int(characterId)), e)
}
