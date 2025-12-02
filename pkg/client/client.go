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

	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// Sets the URI for service account authentication using OAuth for Itential
	// Platform servers.
	tokenUrl = "oauth/token"

	// Sets the URI for authenticating to the Itential Platform server using
	// basic authentication
	authUrl = "/login"

	// Sets the value for the scheme (protocol) to use when constructing the
	// full URL to the server to `http`
	ProtocolHttp = "http"

	// Sets the value for the scheme (protocol) to use when constructing the
	// full URL to the server to `https`
	ProtocolHttps = "https"
)

// HttpClient is the default (and only) client for connecting to Itential
// Platform.  It is defined by the set configuation values and implements the
// Client interface
type HttpClient struct {
	// Host is the hostname or IP address of the Itential Platform server to
	// connect to when the application runs
	Host string

	// Port is the port to use when conneting to the server.  If the port value
	// is set to 0, it will be automatically determined by the value of UseTls.
	// When UseTls is set to true, Port will be set to 443 and when UseTls is
	// set to false, Port will be set to 80
	Port int

	// UseTls enables or disables the use of TLS when connecting to the
	// Itential Platform server.  This value is also used to auto determine the
	// value for Port if Port is set to 0 (default).
	UseTls bool

	// Verify enables or disables certificate verification when connecting
	// to the Itential Plaform server.  This field is only used when UseTls is
	// set to true.
	Verify bool

	// Username is used as the username when authenticating to the Itential
	// Platform server using basic authorization.
	Username string

	// Password is the password to use when authenticating to the Itential
	// Platform server using basic authorization
	Password string

	// ClientId is the client id to use when authenticating to the Itential
	// Platform server using OAuth
	ClientId string

	// ClientSecret is the clietn secret to use when authenticating to the
	// Itential Platform server using OAuth
	ClientSecret string

	// jar is the cookie jar associated with the http session
	jar *Jar

	// authenicated is a flag that is set when the client successfully
	// authenicates to an Itential Platform server
	authenticated bool

	// context is the context for the client.
	context context.Context
}

// New will create a new instance of HttpClient and return it to the calling
// function.  the `ctx` argument will set the context for the HTTP session for
// the applicaiton.  The `cfg` argument provides the configuration profile
// defined from the application configuration.   The configuration profile
// provides the discrete settings for the HTTP client.
func New(ctx context.Context, cfg *config.Profile) *HttpClient {
	logger.Info("Creating new http client")

	// create the client object and return it to the calling function.
	return &HttpClient{
		Host:          cfg.Host,
		Port:          cfg.Port,
		Username:      cfg.Username,
		Password:      cfg.Password,
		UseTls:        cfg.UseTLS,
		Verify:        cfg.Verify,
		ClientId:      cfg.ClientID,
		ClientSecret:  cfg.ClientSecret,
		context:       ctx,
		jar:           NewJar(),
		authenticated: false,
	}
}

// IsAuthenticated returns whether or not the HTTP client has successfully
// authenticated to the Itential Platform server
func (c *HttpClient) IsAuthenticated() bool {
	return c.authenticated
}

// send will send the request to the Itential Platform server and return the
// results and error to the calling function.  The `method` argument specifies
// the HTTP method to use when sending the request.  The `request` argument
// provides the contents of the request to send to the server.
//
// The function returns a response object that includes details about the
// response received from the server.   If there was an error in the request or
// response, it is returned to the calling function.
func (c *HttpClient) send(method string, request *Request) (*Response, error) {
	logger.Trace()

	var scheme string = c.setScheme()
	var remotePort int = c.setPort()
	var remoteHost string = fmt.Sprintf("%s:%v", c.Host, remotePort)

	// contruct the full URL object that is used to send the request
	u := c.newUrl(scheme, remoteHost, request.Path, request.Params)

	logger.Info("%s %s", method, u.String())

	if !request.NoLog {
		debugOutput := string(request.Body)
		if debugOutput != "" {
			logger.Debug("%s", string(request.Body))
		} else {
			logger.Debug("Request body is empty")
		}
	} else {
		logger.Debug("request body is omitted due to the use of NoLog")
	}

	client := &http.Client{Jar: c.jar}

	// Disable certificate verification when UseTls is true and Verify is
	// false.  This is inherently insecure
	if c.UseTls && !c.Verify {
		logger.Debug("Disabling client certificate verification")
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// attempt to authenticate to the server using oauth
	if c.ClientId != "" && c.ClientSecret != "" {
		logger.Debug("attempting to authenticate using client id")
		httpClient, err := c.authenticateUsingOAuth(client, scheme, remoteHost)
		if err != nil {
			return nil, err
		}
		client = httpClient
	}

	// create the actual http request based on method, url and body.
	req, err := c.newHttpRequest(method, u.String(), request.Body)
	if err != nil {
		return nil, err
	}

	// send the request to server and handle the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("HTTP response is %s", resp.Status)

	return c.newResponse(resp, method, u.String())
}

// newResponse creates a response object to return to the calling function
// based on the HTTP response.  The `r` argument is the HTTP response
// object returned from the server.  The `method` is the HTTP method that use
// used to send the request.  The `u` argument is the full URL path.  This
// function will return a client Response object or an error.
func (c *HttpClient) newResponse(r *http.Response, method, u string) (*Response, error) {
	logger.Trace()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err, "failed to read the response body")
		return nil, err
	}

	headers := make(map[string]string)

	for key, value := range r.Header {
		headers[key] = value[0]
	}

	return &Response{
		Method:     method,
		Url:        u,
		Body:       b,
		StatusCode: r.StatusCode,
		Status:     r.Status,
		Headers:    headers,
	}, nil
}

// setScheme sets the HTTP scheme (protocol) to use when constructing the full
// URL.  This function returns either ProtocolHttp or ProtocolHttps depending
// on the value of UseTls
func (c *HttpClient) setScheme() string {
	logger.Trace()

	if c.UseTls {
		return ProtocolHttps
	} else {
		return ProtocolHttp
	}
}

// setPort will determine the port value to use when connecting to the Itential
// Server.  If the port value is 0, it will be automatically determined based
// on the value of UseTls.  When UseTls is true, the port value will be set to
// 443 and when UseTls is falase, the port value will be set to 80.
func (c *HttpClient) setPort() int {
	logger.Trace()

	var port int = c.Port

	if port == 0 && c.UseTls {
		port = 443
	} else if port == 0 && !c.UseTls {
		port = 80
	} else if port == 0 {
		logger.Fatal(fmt.Errorf("could not determine the value for port"), "")
	}

	return port
}

// newUrl will construct a new URL with the provided arguments.   The `scheme`
// argument is a string that should be either `http` or `https`.  The `host`
// argument is the full hostname or address to the server and may or may not
// include the port value.  the `path` argument is the URI to the resource.
//
// the `params` argument will construct a valid and encoded query string to
// append to the URI for the request and add that to the URL object.
//
// This function returns a URL object.
func (c *HttpClient) newUrl(scheme, host, path string, params map[string]string) url.URL {
	logger.Trace()

	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	q := u.Query()

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()

	return u
}

// newHttpRequest creates a new HTTP request object.  This is the actual
// request object that will be used to send to the server.  The `method`
// argument defines the HTTP method to use in the request.  The `u` argumetn
// deifnes the full URL string to send the request to.  The `body` argument
// defines the actual body to incude in the request.
//
// This function will also update the request with the set of default headers
// designed to work with Itential Platform.
func (c *HttpClient) newHttpRequest(method, u string, body []byte) (*http.Request, error) {
	logger.Trace()

	// create the http request
	req, err := http.NewRequestWithContext(
		c.context,
		method,
		u,
		bytes.NewBuffer(body),
	)

	if err != nil {
		return nil, err
	}

	// update the request headers with standard values
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(string(body))))

	return req, nil
}

// authenicateUsingBasicAuth is responsible for authenticating to the Itential
// Platformserver when basic authentication is used.  The client decides to
// use basic authentication with there is a username and password configured
// and there is not a clientid and clientsecret configured.  This function
// will set the value of authenticate to true if the client successfully
// authenticates.
func (c *HttpClient) authenticateUsingBasicAuth() {
	logger.Trace()

	// Construct the body use to authenticate the session.  Itential Platform
	// does not use a standard basic authentication mechanism, rather it
	// expects the client to send the the authentication credentials in the
	// body of a POST request to the server.
	body := map[string]interface{}{
		"user": map[string]interface{}{
			"username": c.Username,
			"password": c.Password,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		logger.Fatal(err, "error attempting to marshal authentication credentials")
	}

	// Be sure to always set WithNoLog(true) for this call to prevent the
	// authentication credentials from leaking out when `--verbose` is used.
	req := NewRequest(authUrl, WithBody(b), WithNoLog(true))

	res, err := c.send("POST", req)
	if err != nil {
		logger.Fatal(
			fmt.Errorf("error sending POST request to the server for authentication"),
			"",
		)
	}

	if res.StatusCode != http.StatusOK {
		logger.Fatal(
			fmt.Errorf("http returned status code `%v` while attempting to authenticate", res.StatusCode),
			"",
		)
	}

	c.authenticated = true
}

// authenticateUsingOAuth will attempt to authenticate to the Itential Platform
// server using the configured ClientId and ClientSecret.  If the
// authentication is successful, this function will return a http.Client object
// that can be used.  If there is an error, the error is returned.
func (c *HttpClient) authenticateUsingOAuth(httpClient *http.Client, scheme, remoteHost string) (*http.Client, error) {
	logger.Trace()

	if c.ClientId == "" || c.ClientSecret == "" {
		return nil, fmt.Errorf("missing client_id or client_secret, authentication failed")
	}

	cfg := &clientcredentials.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{},
		TokenURL:     fmt.Sprintf("%s://%s/%s", scheme, remoteHost, tokenUrl),
		AuthStyle:    1,
	}

	return cfg.Client(context.WithValue(
		c.context,
		oauth2.HTTPClient,
		httpClient,
	)), nil
}

func (c *HttpClient) request(method string, request *Request) (*Response, error) {
	if !c.authenticated && (c.ClientId == "" && c.ClientSecret == "") {
		c.authenticateUsingBasicAuth()
	}
	return c.send(method, request)
}

// Post implements the client interface Post function
func (c *HttpClient) Post(req *Request) (*Response, error) {
	return c.request(http.MethodPost, req)
}

// Get implements the client interface Get function
func (c *HttpClient) Get(req *Request) (*Response, error) {
	return c.request(http.MethodGet, req)
}

// Put implements the client interface Put function
func (c *HttpClient) Put(req *Request) (*Response, error) {
	return c.request(http.MethodPut, req)
}

// Delete implements the client interface Delete function
func (c *HttpClient) Delete(req *Request) (*Response, error) {
	return c.request(http.MethodDelete, req)
}

// Patch implements the client interface Patch function
func (c *HttpClient) Patch(req *Request) (*Response, error) {
	return c.request(http.MethodPatch, req)
}

// Trace implements the client interface Trace function
func (c *HttpClient) Trace(req *Request) (*Response, error) {
	return c.request(http.MethodTrace, req)
}
