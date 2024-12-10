package dto

type EnricherArgType string

const (
	URL      EnricherArgType = "URL"
	HASH     EnricherArgType = "HASH"
	FILE     EnricherArgType = "FILE"
	USERNAME EnricherArgType = "USERNAME"
)

type EnricherInputData struct {
	WebhookUri string
	Data       string
	DataType   EnricherArgType `json:"type"`
}

type EnricherResult struct {
	Report map[string]interface{}
	Errors []string
}
