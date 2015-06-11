package constants

const (
	PARAM_STRING = iota
	PARAM_INT
	PARAM_INT64
	PARAM_FLOAT
	PARAM_FLOAT64
	PARAM_BOOL

	IN_PATH = iota
	IN_QUERY
	IN_HEADER
	IN_BODY

	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PUT    = "PUT"
	METHOD_DELETE = "DELETE"

	SCHEME_HTTP  = "http"
	SCHEME_HTTPS = "https"

	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"

	HEADER_Accept                        = "Accept"
)
