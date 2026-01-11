// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

// Response represents an HTTP response received from the Itential Platform.
//
// Response contains the complete HTTP response including status, headers,
// and body. The response body is provided as raw bytes and should be
// unmarshaled by the caller if it contains JSON or other structured data.
//
// Important: HTTP error status codes (4xx, 5xx) are returned as successful
// Response instances, not as errors. Callers must check the StatusCode field
// to determine if the request was successful from an HTTP perspective.
//
// Example:
//
//	resp, err := client.Get(req)
//	if err != nil {
//	    // Network error, timeout, or other client-level failure
//	    return err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//	    // HTTP-level error (404, 500, etc.)
//	    return fmt.Errorf("request failed: %s", resp.Status)
//	}
//
//	var data MyStruct
//	if err := json.Unmarshal(resp.Body, &data); err != nil {
//	    return err
//	}
type Response struct {
	// Method is the HTTP method used for the request (GET, POST, PUT, etc.).
	Method string

	// Url is the complete URL that was requested, including scheme, host,
	// path, and query parameters.
	Url string

	// Status is the HTTP status text (e.g., "200 OK", "404 Not Found").
	// This is the descriptive text associated with the StatusCode.
	Status string

	// StatusCode is the HTTP status code returned by the server (e.g., 200, 404, 500).
	// Callers should check this field to determine request success.
	StatusCode int

	// Headers contains the HTTP response headers as key-value pairs.
	// Note: If a header has multiple values, only the first value is included.
	// Header names are case-sensitive as returned by the server.
	Headers map[string]string

	// Body contains the raw response body as bytes.
	// For JSON responses, use json.Unmarshal to decode the body.
	// For empty responses (e.g., 204 No Content), Body will be an empty slice.
	Body []byte
}
