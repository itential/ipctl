// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

type HttpOption func(c *HttpClient)

func withPort(v int) HttpOption {
	return func(c *HttpClient) {
		c.Port = v
	}
}

func withUseTls(v bool) HttpOption {
	return func(c *HttpClient) {
		c.UseTls = v
	}
}

func withVerify(v bool) HttpOption {
	return func(c *HttpClient) {
		c.Verify = v
	}
}

func withUsername(v string) HttpOption {
	return func(c *HttpClient) {
		c.Username = v
	}
}

func withPassword(v string) HttpOption {
	return func(c *HttpClient) {
		c.Password = v
	}
}

func withClientID(v string) HttpOption {
	return func(c *HttpClient) {
		c.ClientId = v
	}
}

func withClientSecret(v string) HttpOption {
	return func(c *HttpClient) {
		c.ClientSecret = v
	}
}
