package member

type Model struct {
	id          uint32
	characterId uint32
	worldId     byte
	channelId   byte
	online      bool
}

func (m Model) CharacterId() uint32 {
	return m.characterId
}

func (m *Model) UpdateStatus(worldId byte, channelId byte, online bool) *Model {
	return &Model{
		characterId: m.characterId,
		worldId:     worldId,
		channelId:   channelId,
		online:      online,
	}
}

func (m Model) WorldId() byte {
	return m.worldId
}

func (m Model) ChannelId() byte {
	return m.channelId
}

func (m Model) Online() bool {
	return m.online
}

func (m Model) Id() uint32 {
	return m.id
}

func NewModel(id uint32, characterId uint32, worldId byte, channelId byte) *Model {
	return &Model{
		id:          id,
		characterId: characterId,
		worldId:     worldId,
		channelId:   channelId,
		online:      true,
	}
}
