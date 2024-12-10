package dto

type EnricherConfigArgType string

const (
	StringArg EnricherConfigArgType = "string"
	IntArg    EnricherConfigArgType = "int"
	BoolArg   EnricherConfigArgType = "bool"
	FloatArg  EnricherConfigArgType = "float"
)

type EnricherConfigArg struct {
	Name         string
	Type         EnricherConfigArgType
	Required     bool
	DefaultValue any `json:"defaultValue,omitempty"`
}

type Enricher struct {
	Enabled        bool
	Name           string
	ExecutablePath string
	Timeout        int64
	AllowedTypes   []EnricherArgType
	Author         string `json:"author,omitempty"`
	Source         string `json:"source,omitempty"`
	Description    string `json:"description,omitempty"`
	ConfigArgs     []EnricherConfigArg
}
