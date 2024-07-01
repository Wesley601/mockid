package entities

type Mappings struct {
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	FileName string
	Index    int
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}
