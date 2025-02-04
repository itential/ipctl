// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

type Response struct {
	Method     string
	Url        string
	Status     string
	StatusCode int
	Headers    map[string]string
	Body       []byte
}
