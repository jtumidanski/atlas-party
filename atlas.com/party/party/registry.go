package party

import (
	"atlas-party/party/member"
	"errors"
	"sync"
)

type registry struct {
	parties        map[uint32]*Model
	characterParty map[uint32]uint32
	lock           sync.RWMutex
}

var once sync.Once
var reg *registry

var runningId = uint32(1000000001)
var runningMemberId = uint32(1000000001)

func GetRegistry() *registry {
	once.Do(func() {
		reg = &registry{
			parties:        make(map[uint32]*Model, 0),
			characterParty: make(map[uint32]uint32),
			lock:           sync.RWMutex{},
		}
	})
	return reg
}

func (r *registry) Create(worldId byte, channelId byte, characterId uint32) *Model {
	r.lock.Lock()
	id := r.getNextId()

	members := []*member.Model{member.NewModel(r.getNextMemberId(), characterId, worldId, channelId)}
	party := &Model{
		id:       id,
		leaderId: characterId,
		members:  members,
	}
	r.parties[id] = party
	r.characterParty[characterId] = id
	r.lock.Unlock()
	return party
}

func (r *registry) Destroy(id uint32) (*Model, error) {
	r.lock.Lock()
	if val, ok := r.parties[id]; ok {
		delete(r.parties, id)
		for _, m := range val.Members() {
			delete(r.characterParty, m.CharacterId())
		}
		r.lock.Unlock()
		return val, nil
	}
	r.lock.Unlock()
	return nil, errors.New("unable to locate party")
}

func (r *registry) Join(id uint32, characterId uint32, worldId byte, channelId byte) (*Model, error) {
	r.lock.Lock()
	if val, ok := r.parties[id]; ok {
		p, err := val.AddMember(r.getNextMemberId(), characterId, worldId, channelId)
		if err != nil {
			return nil, err
		}
		r.parties[id] = p
		r.characterParty[characterId] = id

		r.lock.Unlock()
		return p, nil
	}
	r.lock.Unlock()
	return nil, errors.New("unable to locate party")
}

func (r *registry) Leave(characterId uint32) (*Model, *Model, error) {
	r.lock.Lock()
	if id, ok := r.characterParty[characterId]; ok {
		if val, ok := r.parties[id]; ok {
			if val.LeaderId() == characterId {
				delete(r.parties, id)
				for _, m := range val.Members() {
					delete(r.characterParty, m.CharacterId())
				}
				r.lock.Unlock()
				return val, nil, nil
			} else {
				p, err := val.RemoveMember(characterId)
				if err != nil {
					return nil, nil, err
				}
				if len(val.Members()) == 0 {
					delete(r.parties, val.Id())
				} else {
					r.parties[id] = p
				}
				delete(r.characterParty, characterId)
				r.lock.Unlock()
				return val, p, nil
			}
		}
		r.lock.Unlock()
		return nil, nil, errors.New("unable to party")
	}
	r.lock.Unlock()
	return nil, nil, errors.New("unable to party for character")
}

func (r *registry) UpdateStatus(characterId uint32, worldId byte, channelId byte, online bool) (*Model, error) {
	r.lock.Lock()
	if id, ok := r.characterParty[characterId]; ok {
		if val, ok := r.parties[id]; ok {
			p, err := val.UpdateMemberStatus(characterId, worldId, channelId, online)
			if err != nil {
				return nil, err
			}
			r.lock.Unlock()
			return p, nil
		}
		r.lock.Unlock()
		return nil, errors.New("unable to locate party")
	}
	r.lock.Unlock()
	return nil, errors.New("unable to party for character")
}

func (r *registry) PromoteNewLeader(id uint32, characterId uint32) (*Model, error) {
	r.lock.Lock()
	if val, ok := r.parties[id]; ok {
		p, err := val.PromoteLeader(characterId)
		if err != nil {
			return nil, err
		}
		r.lock.Unlock()
		return p, nil
	}
	r.lock.Unlock()
	r.lock.Unlock()
	return nil, errors.New("unable to locate party")
}

func (r *registry) Get(id uint32) (Model, error) {
	r.lock.RLock()
	if val, ok := r.parties[id]; ok {
		r.lock.RUnlock()
		return *val, nil
	}
	r.lock.RUnlock()
	return Model{}, errors.New("unable to locate party")
}

func (r *registry) GetAll() []Model {
	r.lock.RLock()
	result := make([]Model, 0)
	for _, p := range r.parties {
		result = append(result, *p)
	}
	r.lock.RUnlock()
	return result
}

func (r *registry) GetForMember(characterId uint32) (Model, error) {
	r.lock.RLock()
	if pid, ok := r.characterParty[characterId]; ok {
		if val, ok := r.parties[pid]; ok {
			r.lock.RUnlock()
			return *val, nil
		}
	}
	r.lock.RUnlock()
	return Model{}, errors.New("unable to locate party for member")
}

func (r *registry) getNextId() uint32 {
	ids := existingIds(r.parties)

	var currentId = runningId
	for contains(ids, currentId) {
		currentId = currentId + 1
		if currentId > 2000000000 {
			currentId = 1000000001
		}
		runningId = currentId
	}
	return runningId
}

func (r *registry) getNextMemberId() uint32 {
	ids := existingMemberIds(r.parties)

	var currentId = runningMemberId
	for contains(ids, currentId) {
		currentId = currentId + 1
		if currentId > 2000000000 {
			currentId = 1000000001
		}
		runningMemberId = currentId
	}
	return runningMemberId
}

func existingIds(existing map[uint32]*Model) []uint32 {
	var ids []uint32
	for _, x := range existing {
		ids = append(ids, x.Id())
	}
	return ids
}

func existingMemberIds(existing map[uint32]*Model) []uint32 {
	var ids []uint32
	for _, x := range existing {
		for _, y := range x.Members() {
			ids = append(ids, y.Id())
		}
	}
	return ids
}

func contains(ids []uint32, id uint32) bool {
	for _, element := range ids {
		if element == id {
			return true
		}
	}
	return false
}
