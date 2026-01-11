// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"net/http"
	"net/url"
	"sync"
)

// Jar implements the http.CookieJar interface with thread-safe cookie storage.
//
// Jar stores cookies on a per-host basis and provides concurrent access
// through read-write mutex synchronization. It implements a simple cookie
// storage mechanism that stores all cookies for a given host together.
//
// This implementation is thread-safe and can be safely used from multiple
// goroutines. The Cookies method uses a read lock (RLock) to allow concurrent
// reads, while SetCookies uses an exclusive lock for writes.
//
// Note: This implementation does not perform cookie expiration checking,
// path/domain matching, or other RFC 6265 compliance features. It provides
// basic cookie storage suitable for session management with Itential Platform.
type Jar struct {
	lk      sync.RWMutex
	cookies map[string][]*http.Cookie
}

// NewJar creates and returns a new thread-safe cookie jar with an empty
// cookie store.
//
// The returned Jar implements the http.CookieJar interface and can be used
// with http.Client instances. It is safe for concurrent use from multiple
// goroutines.
//
// Example:
//
//	jar := NewJar()
//	client := &http.Client{Jar: jar}
func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies stores cookies for the given URL's host.
//
// This method implements the http.CookieJar interface. It stores all cookies
// for a given host, replacing any previously stored cookies for that host.
// The URL's host (including port if present) is used as the storage key.
//
// This method is thread-safe and uses an exclusive lock to prevent concurrent
// modifications.
//
// Parameters:
//   - u: The URL for which cookies should be stored (only u.Host is used)
//   - cookies: The slice of cookies to store for the host
//
// Example:
//
//	u, _ := url.Parse("https://example.com:8080/api")
//	cookies := []*http.Cookie{{Name: "session", Value: "abc123"}}
//	jar.SetCookies(u, cookies)
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
//
// This method implements the http.CookieJar interface. It retrieves all
// cookies stored for the URL's host. The returned slice may be empty if
// no cookies exist for the host.
//
// This method is thread-safe and uses a read lock to allow concurrent
// reads while preventing data races with concurrent writes.
//
// Note: This implementation does not perform cookie expiration checking
// or path/domain matching as specified in RFC 6265. All stored cookies
// for the host are returned regardless of the URL path.
//
// Parameters:
//   - u: The URL for which to retrieve cookies (only u.Host is used)
//
// Returns:
//   - A slice of cookies for the host, or an empty slice if none exist
//
// Example:
//
//	u, _ := url.Parse("https://example.com/api/users")
//	cookies := jar.Cookies(u)
//	for _, cookie := range cookies {
//	    fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
//	}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	jar.lk.RLock()
	defer jar.lk.RUnlock()
	return jar.cookies[u.Host]
}
