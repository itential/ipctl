// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/itential/ipctl/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	tokenUrl = "oauth/token"
)

type HttpClient struct {
	Host          string
	Port          int
	UseTls        bool
	Verify        bool
	Username      string
	Password      string
	ClientId      string
	ClientSecret  string
	jar           *Jar
	authenticated bool
	context       context.Context
}

func NewHttpClient(ctx context.Context, host string, opts ...HttpOption) *HttpClient {
	logger.Info("Creating new http client")

	httpClient := HttpClient{context: ctx, Host: host, jar: NewJar()}

	for _, opt := range opts {
		opt(&httpClient)
	}

	return &httpClient
}

func (c *HttpClient) authenticate() {
	logger.Trace()

	creds := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
	}

	user := map[string]interface{}{"user": creds}

	b, err := json.Marshal(user)
	if err != nil {
		logger.Fatal(err, "could not marshal credentials")
	}

	req := NewRequest("/login", WithBody(b))

	res, err := c.send("POST", req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	c.authenticated = true
}

func (c *HttpClient) send(method string, request *Request) (*Response, error) {
	logger.Trace()

	var scheme string

	if c.UseTls {
		scheme = "https"
	} else {
		scheme = "http"
	}

	var remoteHost string = c.Host
	if c.Port != 0 {
		remoteHost += ":" + strconv.Itoa(c.Port)
	}

	u := url.URL{
		Scheme: scheme,
		Host:   remoteHost,
		Path:   request.Path,
	}

	q := u.Query()

	if len(request.Params) > 0 {
		for key, value := range request.Params {
			q.Set(key, value)
		}
	}

	u.RawQuery = q.Encode()

	logger.Info("%s %s", method, u.String())

	debugOutput := string(request.Body)

	if debugOutput != "" {
		// FIXME (privateip) this is a quick hack to prevent the logger from
		// showing the username and password used to authenticate to the
		// server.  A better mechanism should be implemented
		if request.Path != "/login" {
			logger.Debug(string(request.Body))
		}
	} else {
		logger.Debug("Request body is empty")
	}

	client := &http.Client{Jar: c.jar}

	if c.UseTls && !c.Verify {
		logger.Debug("Disabling client certificate verification")
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if !c.authenticated && c.ClientId != "" && c.ClientSecret != "" {
		logger.Debug("attempting to authenticate using client id")

		cfg := &clientcredentials.Config{
			ClientID:     c.ClientId,
			ClientSecret: c.ClientSecret,
			Scopes:       []string{},
			TokenURL:     fmt.Sprintf("%s://%s/%s", scheme, remoteHost, tokenUrl),
			AuthStyle:    1,
		}

		client = cfg.Client(context.WithValue(
			c.context,
			oauth2.HTTPClient,
			client,
		))
	}

	req, _ := http.NewRequestWithContext(
		c.context,
		method,
		u.String(),
		bytes.NewBuffer(request.Body),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(string(request.Body))))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("HTTP response is %s", resp.Status)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "read body failed")
		return nil, err
	}

	headers := make(map[string]interface{})

	for key, value := range resp.Header {
		headers[key] = value[0]
	}

	response := &Response{
		Method:     method,
		Url:        u.String(),
		Body:       b,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    map[string]string{},
	}

	for key, value := range resp.Header {
		response.Headers[key] = value[0]
	}

	return response, nil
}

func (c *HttpClient) request(method string, request *Request) (*Response, error) {
	if !c.authenticated && (c.ClientId == "" && c.ClientSecret == "") {
		c.authenticate()
	}
	return c.send(method, request)
}

func (c *HttpClient) Post(req *Request) (*Response, error) {
	return c.request(http.MethodPost, req)
}

func (c *HttpClient) Get(req *Request) (*Response, error) {
	return c.request(http.MethodGet, req)
}

func (c *HttpClient) Put(req *Request) (*Response, error) {
	return c.request(http.MethodPut, req)
}

func (c *HttpClient) Delete(req *Request) (*Response, error) {
	return c.request(http.MethodDelete, req)
}

func (c *HttpClient) Patch(req *Request) (*Response, error) {
	return c.request(http.MethodPatch, req)
}

func (c *HttpClient) Trace(req *Request) (*Response, error) {
	return c.request(http.MethodTrace, req)
}
