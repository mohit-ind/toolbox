package constants

// BaseUrls & paths needed to call api's connectors
const (
	// BaseURLTest points to the combined API Gateway's Test stage
	BaseURLTest = "https://3dr3jb9xh0.execute-api.eu-central-1.amazonaws.com/test"

	// BaseURLStaging points to the combined API Gateway's Staging stage
	BaseURLStaging = "https://3dr3jb9xh0.execute-api.eu-central-1.amazonaws.com/staging"

	// BaseURLOldStaging points to the old (manual) Staging's API Gateway
	BaseURLOldStaging = "https://48oy28ta73.execute-api.eu-central-1.amazonaws.com/stage"

	// BaseURLStaging points to the combined API Gateway's Acceptance stage
	BaseURLAcceptance = "https://3dr3jb9xh0.execute-api.eu-central-1.amazonaws.com/acceptance"

	// BaseURLProduction points to the Production API Gateway
	BaseURLProduction = "https://api.appventurez.nl"
)

const (
	// DecodeTokenPath to UMS's decode Cognito JWT controller
	DecodeTokenPath = "/decode-token"

	// LabelInfoPath to UMS's GetLabelInfo controller
	LabelInfoPath = "/label/info"
)
