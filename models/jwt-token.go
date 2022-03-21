package models

// DecodeTokenRequest ... here we have to provide JWT token like this "Bearer  eyJraWQiOi....."
type DecodeTokenRequest struct {
	// Token is the JWT Token string
	Token string `json:"token"`
}

// DecodeTokenResponse ... have fields comes in response while we call Decode Token API which fills when we call 'Do' func in connector and decode the response in it.
// TODO: Add meaningful dockstring to exported field names
type DecodeTokenResponse struct {
	Email              string   `json:"user_email"`
	UserID             string   `json:"user_id"`
	OrgID              string   `json:"org_id"`
	UserAutobinckID    string   `json:"user_autobinckID"`
	OrgAutobinckID     string   `json:"org_autobinck_id"`
	CustomerIdentifier string   `json:"customer_id"`
	TravelerIdentifier string   `json:"traveler_id"`
	UserRoles          []string `json:"user_roles"`
	UMSGoldenSource    bool     `json:"ums_golden_source"`
	ForcedLogout       bool     `json:"forced_logout"`
}
