package models

type ApiVersion struct {
	Version      string   `json:"version"`
	SupportedCRs []string `json:"supportedCRs"`
	BuildTime    string   `json:"buildTime"`
	CommitHash   string   `json:"commitHash"`
	ServiceName  string   `json:"serviceName"`
}
