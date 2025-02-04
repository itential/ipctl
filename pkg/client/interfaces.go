// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

type Client interface {
	Http() Http
	Profile() map[string]interface{}
}

type Http interface {
	Get(*Request) (*Response, error)
	Post(*Request) (*Response, error)
	Put(*Request) (*Response, error)
	Delete(*Request) (*Response, error)
	Patch(*Request) (*Response, error)
	Trace(*Request) (*Response, error)
}
