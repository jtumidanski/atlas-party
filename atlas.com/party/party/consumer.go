package party

import (
	"atlas-party/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNamePartyCreate        = "party_create_command"
	consumerNamePartyExpel         = "party_expel_command"
	consumerNamePartyJoin          = "party_join_command"
	consumerNamePartyLeave         = "party_leave_command"
	consumerNamePartyPromoteLeader = "party_promote_leader_command"
	topicTokenPartyCreate          = "TOPIC_PARTY_CREATE"
	topicTokenPartyExpel           = "TOPIC_PARTY_EXPEL"
	topicTokenPartyJoin            = "TOPIC_PARTY_JOIN"
	topicTokenPartyLeave           = "TOPIC_PARTY_LEAVE"
	topicTokenPromoteLeader        = "TOPIC_PARTY_PROMOTE_LEADER"
)

func CreateConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[partyCreateCommand](consumerNamePartyCreate, topicTokenPartyCreate, groupId, handleCreate())
}

func handleCreate() kafka.HandlerFunc[partyCreateCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command partyCreateCommand) {
		Create(l, span)(command.CharacterId, command.WorldId, command.ChannelId)
	}
}

func ExpelConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[partyExpelCommand](consumerNamePartyExpel, topicTokenPartyExpel, groupId, handleExpel())
}

type partyExpelCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	PartyId     uint32 `json:"party_id"`
}

func handleExpel() kafka.HandlerFunc[partyExpelCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command partyExpelCommand) {
		Expel(l, span)(command.WorldId, command.ChannelId, command.CharacterId)
	}
}

func JoinConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[partyJoinCommand](consumerNamePartyJoin, topicTokenPartyJoin, groupId, handleJoin())
}

func handleJoin() kafka.HandlerFunc[partyJoinCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command partyJoinCommand) {
		Join(l, span)(command.WorldId, command.ChannelId, command.PartyId, command.CharacterId)
	}
}

func LeaveConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[partyLeaveCommand](consumerNamePartyLeave, topicTokenPartyLeave, groupId, handleLeave())
}

func handleLeave() kafka.HandlerFunc[partyLeaveCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command partyLeaveCommand) {
		Leave(l, span)(command.WorldId, command.ChannelId, command.CharacterId)
	}
}

func PromoteLeaderConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[partyPromoteLeaderCommand](consumerNamePartyPromoteLeader, topicTokenPromoteLeader, groupId, handlePromoteLeader())
}

type partyPromoteLeaderCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	PartyId     uint32 `json:"party_id"`
}

func handlePromoteLeader() kafka.HandlerFunc[partyPromoteLeaderCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command partyPromoteLeaderCommand) {
		PromoteNewLeader(l, span)(command.WorldId, command.ChannelId, command.PartyId, command.CharacterId)
	}
}
