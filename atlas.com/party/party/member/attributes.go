package member

type Attributes struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	Online      bool   `json:"online"`
}

func MakeAttribute(m *Model) Attributes {
	return Attributes{
		WorldId:     m.WorldId(),
		ChannelId:   m.ChannelId(),
		CharacterId: m.CharacterId(),
		Online:      m.Online(),
	}
}

func MakeAttributes(members []*Model) []Attributes {
	result := make([]Attributes, 0)
	for _, m := range members {
		result = append(result, MakeAttribute(m))
	}
	return result
}