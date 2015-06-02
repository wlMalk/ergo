package ergo

const (
	IN_PATH = iota
	IN_QUERY
	IN_HEADER
	IN_BODY

	METHOD_ANY    = ""
	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PUT    = "PUT"
	METHOD_DELETE = "DELETE"

	SCHEME_HTTP  = "http"
	SCHEME_HTTPS = "https"

	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"

)
