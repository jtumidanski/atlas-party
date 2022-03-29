package party

import (
	"atlas-party/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type partyCreateCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func emitCreateParty(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_CREATE")
	return func(worldId byte, channelId byte, characterId uint32) {
		e := &partyCreateCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type partyJoinCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	PartyId     uint32 `json:"party_id"`
}

func emitJoinParty(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_JOIN")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		e := &partyJoinCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			PartyId:     partyId,
			CharacterId: characterId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type partyLeaveCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func emitLeaveParty(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_LEAVE")
	return func(worldId byte, channelId byte, characterId uint32) {
		e := &partyLeaveCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type partyMemberStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func emitPartyMemberLogin(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "LOGIN")
	}
}

func emitPartyMemberLogout(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "LOGOUT")
	}
}

func emitPartyMemberJoin(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "JOINED")
	}
}

func emitPartyMemberLeave(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "LEFT")
	}
}

func emitPartyMemberExpelled(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "EXPELLED")
	}
}

func emitPartyMemberDisbanded(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		emitPartyMemberStatus(producer, worldId, channelId, partyId, characterId, "DISBANDED")
	}
}

func emitPartyMemberPromoted(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_MEMBER_STATUS")
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
	producer(kafka.CreateKey(int(characterId)), e)
}

type partyStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func emitPartyCreated(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_STATUS")
	return func(worldId byte, partyId uint32, characterId uint32) {
		emitPartyStatus(producer, worldId, partyId, characterId, "CREATED")
	}
}

func emitPartyDisbanded(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, partyId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PARTY_STATUS")
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
	producer(kafka.CreateKey(int(partyId)), e)
}
