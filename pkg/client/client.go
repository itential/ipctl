// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"context"

	"github.com/itential/ipctl/pkg/config"
)

type IapClient struct {
	http *HttpClient
}

func New(ctx context.Context, cfg *config.Profile) *IapClient {
	return &IapClient{
		http: NewHttpClient(
			ctx,
			cfg.Host,
			withPort(cfg.Port),
			withUsername(cfg.Username),
			withPassword(cfg.Password),
			withUseTls(cfg.UseTLS),
			withVerify(cfg.Verify),
			withClientID(cfg.ClientID),
			withClientSecret(cfg.ClientSecret),
		),
	}
}

func (c *IapClient) Http() Http {
	return c.http
}

func (c *IapClient) Profile() map[string]interface{} {
	var port int

	if c.http.UseTls {
		port = 443
	} else if !c.http.UseTls {
		port = 80
	} else {
		port = c.http.Port
	}

	return map[string]interface{}{
		"host":          c.http.Host,
		"port":          port,
		"use_tls":       c.http.UseTls,
		"verify":        c.http.Verify,
		"username":      c.http.Username,
		"password":      "********",
		"client_id":     c.http.ClientId,
		"client_secret": "********",
	}
}

func (c *IapClient) IsAuthenticated() bool {
	return c.http.authenticated
}
