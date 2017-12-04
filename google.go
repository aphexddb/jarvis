package jarvis

// GoogleHomeWebhookRequest is a request from Google Home
type GoogleHomeWebhookRequest struct {
	ResponseID  string `json:"responseId"`
	QueryResult struct {
		QueryText                string            `json:"queryText"`
		Action                   string            `json:"action"`
		Parameters               map[string]string `json:"parameters"`
		AllRequiredParamsPresent bool              `json:"allRequiredParamsPresent"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text"`
			} `json:"text"`
		} `json:"fulfillmentMessages"`
		Intent struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"intent"`
		IntentDetectionConfidence float64 `json:"intentDetectionConfidence"`
		DiagnosticInfo            struct {
		} `json:"diagnosticInfo"`
		LanguageCode string `json:"languageCode"`
	} `json:"queryResult"`
	OriginalDetectIntentRequest struct {
		Payload struct {
		} `json:"payload"`
	} `json:"originalDetectIntentRequest"`
	Session string `json:"session"`
}
