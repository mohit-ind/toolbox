package models

// LabelBasedInfoResponse - It have all label fields which are essential to provide info about any label, with error which comes due to bad connection or any other reasons like invalid api-key.
// It also have customer level info based on that api key.
// TODO: Add meaningful godoc to exported fields
type LabelBasedInfoResponse struct {
	Source           string            `json:"source"`
	LabelAutobinckID string            `json:"label_autobinck_id"`
	LabelKey         string            `json:"label_key"`
	CustomerInfo     *CustomerInfo     `json:"customer_info,omitempty"`
	Error            ErrorResponseBody `json:"error"`
}

// CustomerInfo - is a general structure which have all customer data comes during calling connector to get label-info (Only if you provide 'customer-identifier' in query-param).
// TODO: Add meaningful godoc to exported fields
type CustomerInfo struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	AutobinckID      string `json:"autobinck_id"`
	LabelAutobinckID string `json:"label_autobinck_id"`
	ExternalID       string `json:"external_id"`
	KvkNumber        string `json:"kvk_number"`
	VatNumber        string `json:"vat_number"`
	MobinckID        string `json:"mobinck_id"`
	UMSGoldenSource  bool   `json:"ums_golden_source"`
}

// ErrorResponseBody - Basic error response body format according to appventurez standard.
// TODO: Add meaningful godoc to exported fields
type ErrorResponseBody struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	MessageCode string `json:"messageCode"`
}
