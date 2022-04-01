package party

import (
	"atlas-party/model"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func ByIdModelProvider(_ logrus.FieldLogger, _ opentracing.Span) func(partyId uint32) model.Provider[Model] {
	return func(partyId uint32) model.Provider[Model] {
		return func() (Model, error) {
			return GetRegistry().Get(partyId)
		}
	}
}

func ByMemberModelProvider(_ logrus.FieldLogger, _ opentracing.Span) func(characterId uint32) model.Provider[Model] {
	return func(characterId uint32) model.Provider[Model] {
		return func() (Model, error) {
			return GetRegistry().GetForMember(characterId)
		}
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(partyId uint32) (Model, error) {
	return func(partyId uint32) (Model, error) {
		return ByIdModelProvider(l, span)(partyId)()
	}
}

func GetByMember(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		return ByMemberModelProvider(l, span)(characterId)()
	}
}

func GetAll() []Model {
	return GetRegistry().GetAll()
}

func Create(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, worldId byte, channelId byte) {
	return func(characterId uint32, worldId byte, channelId byte) {
		p := GetRegistry().Create(worldId, channelId, characterId)
		l.Debugf("Party %d created by character %d in world %d.", p.Id(), p.LeaderId(), worldId)
		emitPartyCreated(l, span)(worldId, p.Id(), characterId)
	}
}

func Leave(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		previous, current, err := GetRegistry().Leave(characterId)
		if err != nil {
			l.WithError(err).Errorf("Character %d was unable to leave their party.", characterId)
			return
		}

		l.Debugf("Character %d left party %d.", characterId, previous.Id())

		if current == nil {
			l.Debugf("As a result, party %d will be disbanded.", previous.Id())
			emitPartyDisbanded(l, span)(previous.Members()[0].WorldId(), previous.Id(), characterId)
			for _, m := range previous.Members() {
				emitPartyMemberDisbanded(l, span)(m.WorldId(), m.ChannelId(), previous.Id(), m.CharacterId())
			}
		} else {
			emitPartyMemberLeave(l, span)(worldId, channelId, previous.Id(), characterId)
		}
	}
}

func Expel(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		previous, _, err := GetRegistry().Leave(characterId)
		if err != nil {
			l.WithError(err).Errorf("Character %d was unable to leave their party, due to expulsion.", characterId)
			return
		}

		l.Debugf("Character %d was expelled from party %d.", characterId, previous.Id())
		emitPartyMemberExpelled(l, span)(worldId, channelId, previous.Id(), characterId)
	}
}

func Join(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		_, err := GetRegistry().Join(partyId, characterId, worldId, channelId)
		if err != nil {
			l.WithError(err).Errorf("Character %d was unable to join party %d.", characterId, partyId)
			return
		}
		emitPartyMemberJoin(l, span)(worldId, channelId, partyId, characterId)
	}
}

func PromoteNewLeader(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) {
		_, err := GetRegistry().PromoteNewLeader(partyId, characterId)
		if err != nil {
			l.WithError(err).Errorf("Character %d was unable to become the new leader of party %d.", characterId, partyId)
			return
		}
		emitPartyMemberPromoted(l, span)(worldId, channelId, partyId, characterId)
	}
}

func MemberLogin(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, worldId byte, channelId byte) {
	return func(characterId uint32, worldId byte, channelId byte) {
		p, err := GetRegistry().UpdateStatus(characterId, worldId, channelId, true)
		if err != nil {
			l.WithError(err).Errorf("Unable to mark character %d as online for party.", characterId)
			return
		}
		emitPartyMemberLogin(l, span)(worldId, channelId, p.Id(), characterId)
	}
}

func MemberLogout(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, worldId byte, channelId byte) {
	return func(characterId uint32, worldId byte, channelId byte) {
		p, err := GetRegistry().UpdateStatus(characterId, 0, 0, false)
		if err != nil {
			l.WithError(err).Errorf("Unable to mark character %d as offline for party.", characterId)
			return
		}
		emitPartyMemberLogout(l, span)(worldId, channelId, p.Id(), characterId)
	}
}
