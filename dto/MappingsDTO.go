package dto

type MappingsDTO struct {
	Index    string              `json:"index"`
	Mappings []MappingsDetailDTO `json:"mappings"`
}

type MappingsDetailDTO struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}
