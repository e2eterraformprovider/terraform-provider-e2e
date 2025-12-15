package constants

// ============================================================================
// API Request Parameter Constants
// ============================================================================
// These constants represent parameter names used in API requests.
//
// SECURITY NOTE: The E2E Networks API requires authentication credentials
// to be passed as query parameters (apikey, project_id, location).
// This is a requirement of the API design and cannot be changed on the client side.
//
// Security implications of query parameter authentication:
// - Query parameters are logged in web server access logs
// - Query parameters may appear in browser history
// - Query parameters may be visible in proxy server logs
// - Query parameters can leak through HTTP Referer headers
//
// Users of this SDK should:
// - Use HTTPS for all API calls (enforced by default)
// - Be aware of logging implications
// - Rotate API keys regularly
// - Use appropriate network security controls

const (
	// QueryParamAPIKey is the query parameter name for the API key.
	// This parameter is required for all API requests.
	QueryParamAPIKey = "apikey"

	// QueryParamProjectID is the query parameter name for the project ID.
	// This parameter is required for all API requests.
	QueryParamProjectID = "project_id"

	// QueryParamLocation is the query parameter name for the location/region.
	// This parameter is required for all API requests.
	QueryParamLocation = "location"

	// QueryParamWorkspaceID is the query parameter name for the workspace ID.
	// This parameter is optional and used for multi-tenant scenarios.
	QueryParamWorkspaceID = "workspace_id"

	// QueryParamTeamID is the query parameter name for the team ID.
	// This parameter is optional and used for team-based scenarios.
	QueryParamTeamID = "team_id"
)
