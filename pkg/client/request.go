// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

// A Request instance is used to send a HTTP request to a remote host.
type Request struct {
	// The request path to be used.  This value is appened to the BaseUrl set
	// in the HTTP client to form the full URL of the request.
	Path string

	// The query parameters used to construct the query string for the request.
	Params map[string]string

	// The HTTP headers to be sent to the remote host
	Headers map[string]string

	// The HTTP body to be send to the remote host
	Body []byte
}

// Defines a new HTTP request object.
func NewRequest(path string, opts ...RequestOption) *Request {
	req := &Request{Path: path}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

type RequestOption func(r *Request)

// Sets the HTTP headers to be sent to the remote host for the specified
// request.
func WithHeaders(v map[string]string) RequestOption {
	return func(r *Request) {
		r.Headers = v
	}
}

// Sets the query string appended to the end of the URL to be sent to the
// remote host for the specified request
func WithParams(v map[string]string) RequestOption {
	return func(r *Request) {
		r.Params = v
	}
}

// Sets the body of the request message to be sent to the remote host for the
// specified request.
func WithBody(v []byte) RequestOption {
	return func(r *Request) {
		r.Body = v
	}
}
