package consumers

import (
	"atlas-party/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	CharacterStatusEvent      = "character_status_event"
	PartyCreateCommand        = "party_create_command"
	PartyExpelCommand         = "party_expel_command"
	PartyJoinCommand          = "party_join_command"
	PartyLeaveCommand         = "party_leave_command"
	PartyPromoteLeaderCommand = "party_promote_leader_command"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHARACTER_STATUS", CharacterStatusEvent, EmptyCharacterStatusEventCreator(), HandleCharacterStatusEvent())
	cec("TOPIC_PARTY_CREATE", PartyCreateCommand, EmptyPartyCreateCommandCreator(), HandlePartyCreateCommand())
	cec("TOPIC_PARTY_EXPEL", PartyExpelCommand, EmptyPartyExpelCommandCreator(), HandlePartyExpelCommand())
	cec("TOPIC_PARTY_JOIN", PartyJoinCommand, EmptyPartyJoinCommandCreator(), HandlePartyJoinCommand())
	cec("TOPIC_PARTY_LEAVE", PartyLeaveCommand, EmptyPartyLeaveCommandCreator(), HandlePartyLeaveCommand())
	cec("TOPIC_PARTY_PROMOTE_LEADER", PartyPromoteLeaderCommand, EmptyPartyPromoteLeaderCommandCreator(), HandlePartyPromoteLeaderCommand())
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, name, topicToken, "Party Orchestration Service", emptyEventCreator, processor)
}
