// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

// Response provides a response object related to a HTTP Request.  The response
// object provides the response from the server with associated metadata.
type Response struct {
	// Method is the method that was used in the original Request object.
	Method string

	// Url is the full URL that was used in to send the request to the sever
	Url string

	// Status is the HTTP describptive text status based on the status code value.
	Status string

	// StatusCode is the HTTP status code returned from the server
	StatusCode int

	// Headers provides the set of headers returned from the server
	Headers map[string]string

	// Body is the actual body returned from the server
	Body []byte
}
