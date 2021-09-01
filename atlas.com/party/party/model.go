package party

import (
	"atlas-party/party/member"
	"errors"
)

type Model struct {
	id       uint32
	leaderId uint32
	members  []*member.Model
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) LeaderId() uint32 {
	return m.leaderId
}

func (m Model) Members() []*member.Model {
	return m.members
}

func (m Model) AddMember(id uint32, characterId uint32, worldId byte, channelId byte) (*Model, error) {
	for _, em := range m.members {
		if em.CharacterId() == characterId {
			return nil, errors.New("character already in party")
		}
	}

	nm := &Model{
		id:       m.Id(),
		leaderId: m.LeaderId(),
		members:  append(m.members, member.NewModel(id, characterId, worldId, channelId)),
	}
	return nm, nil
}

func (m Model) RemoveMember(characterId uint32) (*Model, error) {
	nms := make([]*member.Model, 0)
	found := false
	for _, em := range m.members {
		if em.CharacterId() == characterId {
			found = true
		}
		nms = append(nms, em)
	}
	if !found {
		return nil, errors.New("character not in party")
	}

	leaderId := m.LeaderId()
	if characterId == leaderId {
		if len(nms) == 0 {
			leaderId = 0
		} else {
			leaderId = nms[0].CharacterId()
		}
	}

	nm := &Model{
		id:       m.Id(),
		leaderId: leaderId,
		members:  nms,
	}
	return nm, nil
}

func (m Model) UpdateMemberStatus(characterId uint32, worldId byte, channelId byte, online bool) (*Model, error) {
	nms := make([]*member.Model, 0)
	found := false
	for _, em := range m.members {
		if em.CharacterId() == characterId {
			nms = append(nms, em.UpdateStatus(worldId, channelId, online))
			found = true
		} else {
			nms = append(nms, em)
		}
	}
	if !found {
		return nil, errors.New("character not in party")
	}

	nm := &Model{
		id:       m.Id(),
		leaderId: m.LeaderId(),
		members:  nms,
	}
	return nm, nil
}

func (m Model) PromoteLeader(characterId uint32) (*Model, error) {
	found := false
	for _, em := range m.members {
		if em.CharacterId() == characterId {
			found = true
		}
	}
	if !found {
		return nil, errors.New("character not in party")
	}
	nm := &Model{
		id:       m.Id(),
		leaderId: characterId,
		members:  m.Members(),
	}
	return nm, nil
}
