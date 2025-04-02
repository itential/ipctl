// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

// Client is an interface that represents an Itential Platform client that can
// be used to send and receive messages to and from a server over HTTP.
type Client interface {
	// Get will send a HTTP GET request to the Itental Platform server and
	// return the response
	Get(*Request) (*Response, error)

	// Post will send a HTTP POST request to the Itental Platform server and
	// return the response
	Post(*Request) (*Response, error)

	// Put will send a HTTP PUT request to the Itental Platform server and
	// return the response
	Put(*Request) (*Response, error)

	// Delete will send a HTTP DELETE request to the Itental Platform server and
	// return the response
	Delete(*Request) (*Response, error)

	// Patch will send a HTTP PATCH request to the Itental Platform server and
	// return the response
	Patch(*Request) (*Response, error)

	// Trace will send a HTTP TRACE request to the Itental Platform server and
	// return the response
	Trace(*Request) (*Response, error)
}
