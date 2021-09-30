package structures

// WEBHOOK types
type PostWebhookReqHeaders struct {
	Bearer string `json:"Bearer"`
}

type PostWebhookReq struct {
	Headers PostWebhookReqHeaders `json:"headers"`
	Request interface{}           `json:"request"`
}

type PostWebhookRes struct {
	HasuraUserId  string `json:"X-Hasura-User-Id"`
	HasuraRole    string `json:"X-Hasura-Role"`
	HasuraIsOwner string `json:"X-Hasura-Is-Owner"`
	HasuraCustom  string `json:"X-Hasura-Custom"`
}
