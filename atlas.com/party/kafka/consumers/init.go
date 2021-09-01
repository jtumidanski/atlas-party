package consumers

import (
	"atlas-party/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHARACTER_STATUS", EmptyCharacterStatusEventCreator(), HandleCharacterStatusEvent())
	cec("TOPIC_PARTY_CREATE", EmptyPartyCreateCommandCreator(), HandlePartyCreateCommand())
	cec("TOPIC_PARTY_EXPEL", EmptyPartyExpelCommandCreator(), HandlePartyExpelCommand())
	cec("TOPIC_PARTY_JOIN", EmptyPartyJoinCommandCreator(), HandlePartyJoinCommand())
	cec("TOPIC_PARTY_LEAVE", EmptyPartyLeaveCommandCreator(), HandlePartyLeaveCommand())
	cec("TOPIC_PARTY_PROMOTE_LEADER", EmptyPartyPromoteLeaderCommandCreator(), HandlePartyPromoteLeaderCommand())
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, topicToken, "Party Orchestration Service", emptyEventCreator, processor)
}
