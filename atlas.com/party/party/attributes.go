package party

type Attributes struct {
	LeaderId uint32 `json:"leader_id"`
}

func MakeAttribute(p *Model) Attributes {
	return Attributes{
		LeaderId: p.LeaderId(),
	}
}