package party

type inputDataContainer struct {
	Data inputDataBody `json:"data"`
}

type inputDataBody struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes inputAttributes `json:"attributes"`
}

type inputAttributes struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

type Attributes struct {
	LeaderId uint32 `json:"leader_id"`
}

func MakeAttribute(p Model) Attributes {
	return Attributes{
		LeaderId: p.LeaderId(),
	}
}
