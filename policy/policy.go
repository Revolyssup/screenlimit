package policy

type PolicyRequest struct {
	App    string `json:"app"`
	Action string `json:"action"`
}
