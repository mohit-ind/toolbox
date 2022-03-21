package constants

const (
	// ContextKeyForEmail is used to retrieve the "email" from the context
	ContextKeyForEmail ContextKey = "email"

	// ContextKeyForAutobinckID is used to retrieve the "autobinck_id" from the context
	ContextKeyForAutobinckID ContextKey = "autobinck_id"

	// ContextKeyForUserID is used to retrieve the "user_id" from the context
	ContextKeyForUserID ContextKey = "user_id"

	// ContextKeyForOrgID is used to retrieve the "org_id" from the context
	ContextKeyForOrgID ContextKey = "org_id"

	// ContextKeyForOrgAutobinckID is used to retrieve the "org_autobinck_id" from the context
	ContextKeyForOrgAutobinckID ContextKey = "org_autobinck_id"

	// ContextKeyForCustomerIdentifier is used to retrieve the "customer_identifier" from the context
	ContextKeyForCustomerIdentifier ContextKey = "customer_identifier"

	// ContextKeyForTravelerIdentifier is used to retrieve the "traveler_identifier" from the context
	ContextKeyForTravelerIdentifier ContextKey = "traveler_identifier"

	// ContextKeyForUserRoles is used to retrieve the "user_roles" from the context
	ContextKeyForUserRoles ContextKey = "user_roles"

	// ContextKeyForLabelAutobinckID is used to retrieve the "label_autobinck_id" from the context
	ContextKeyForLabelAutobinckID ContextKey = "label_autobinck_id"

	// ContextKeyForLabelKey is used to retrieve the "label_key" from the context
	ContextKeyForLabelKey ContextKey = "label_key"

	// ContextKeyForSource is used to retrieve the "source" from the context
	ContextKeyForSource ContextKey = "source"

	// ContextKeyForUMSGoldenSource is used to take bool value of UMS golden source for any customer of any traveler
	ContextKeyForUMSGoldenSource ContextKey = "UMS_GOLDEN_SOURCE"

	// ContextKeyForForcedLogout is used to take bool value of forced logout for any traveler
	ContextKeyForForcedLogout ContextKey = "forced_logout"
)

// ContextKey is the key of a context value. The constant context-keys are used to retrieve different values from the http context
type ContextKey string
