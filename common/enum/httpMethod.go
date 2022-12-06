package enum

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "GET"
	HEAD    HttpMethod = "HEAD"
	DELETE  HttpMethod = "DELETE"
	CONNECT HttpMethod = "CONNECT"
	OPTIONS HttpMethod = "OPTIONS"
	TRACE   HttpMethod = "TRACE"
	PATCH   HttpMethod = "PATCH"
)
