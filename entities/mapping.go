package entities

type Mappings struct {
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}
