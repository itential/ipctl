// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestNewJar(t *testing.T) {
	jar := NewJar()

	if jar == nil {
		t.Fatal("Expected NewJar to return a non-nil Jar")
	}

	if jar.cookies == nil {
		t.Error("Expected cookies map to be initialized")
	}

	if len(jar.cookies) != 0 {
		t.Errorf("Expected empty cookies map, got %d entries", len(jar.cookies))
	}
}

func TestJarSetCookies(t *testing.T) {
	jar := NewJar()

	testURL, _ := url.Parse("https://example.com")
	cookies := []*http.Cookie{
		{
			Name:  "session_id",
			Value: "abc123",
		},
		{
			Name:  "user_pref",
			Value: "theme_dark",
		},
	}

	jar.SetCookies(testURL, cookies)

	storedCookies := jar.cookies[testURL.Host]
	if storedCookies == nil {
		t.Fatal("Expected cookies to be stored for host")
	}

	if len(storedCookies) != len(cookies) {
		t.Errorf("Expected %d cookies, got %d", len(cookies), len(storedCookies))
	}

	if !reflect.DeepEqual(storedCookies, cookies) {
		t.Error("Stored cookies do not match original cookies")
	}
}

func TestJarSetCookiesMultipleHosts(t *testing.T) {
	jar := NewJar()

	host1URL, _ := url.Parse("https://example1.com")
	host1Cookies := []*http.Cookie{
		{Name: "host1_session", Value: "session1"},
	}

	host2URL, _ := url.Parse("https://example2.com")
	host2Cookies := []*http.Cookie{
		{Name: "host2_session", Value: "session2"},
	}

	jar.SetCookies(host1URL, host1Cookies)
	jar.SetCookies(host2URL, host2Cookies)

	if len(jar.cookies) != 2 {
		t.Errorf("Expected 2 host entries, got %d", len(jar.cookies))
	}

	stored1 := jar.cookies[host1URL.Host]
	if !reflect.DeepEqual(stored1, host1Cookies) {
		t.Error("Host1 cookies do not match")
	}

	stored2 := jar.cookies[host2URL.Host]
	if !reflect.DeepEqual(stored2, host2Cookies) {
		t.Error("Host2 cookies do not match")
	}
}

func TestJarSetCookiesOverwrite(t *testing.T) {
	jar := NewJar()
	testURL, _ := url.Parse("https://example.com")

	initialCookies := []*http.Cookie{
		{Name: "session", Value: "old_value"},
	}

	newCookies := []*http.Cookie{
		{Name: "session", Value: "new_value"},
		{Name: "token", Value: "auth_token"},
	}

	jar.SetCookies(testURL, initialCookies)
	jar.SetCookies(testURL, newCookies)

	storedCookies := jar.cookies[testURL.Host]
	if len(storedCookies) != len(newCookies) {
		t.Errorf("Expected %d cookies after overwrite, got %d", len(newCookies), len(storedCookies))
	}

	if !reflect.DeepEqual(storedCookies, newCookies) {
		t.Error("Stored cookies should match new cookies after overwrite")
	}
}

func TestJarGetCookies(t *testing.T) {
	jar := NewJar()
	testURL, _ := url.Parse("https://example.com")

	expectedCookies := []*http.Cookie{
		{
			Name:    "session_id",
			Value:   "abc123",
			Path:    "/",
			Domain:  "example.com",
			Expires: time.Now().Add(24 * time.Hour),
		},
		{
			Name:   "preferences",
			Value:  "lang=en",
			Path:   "/",
			Domain: "example.com",
		},
	}

	jar.SetCookies(testURL, expectedCookies)

	retrievedCookies := jar.Cookies(testURL)
	if !reflect.DeepEqual(retrievedCookies, expectedCookies) {
		t.Error("Retrieved cookies do not match expected cookies")
	}
}

func TestJarGetCookiesNoHost(t *testing.T) {
	jar := NewJar()
	testURL, _ := url.Parse("https://nonexistent.com")

	retrievedCookies := jar.Cookies(testURL)
	if retrievedCookies != nil {
		t.Errorf("Expected nil for non-existent host, got %v", retrievedCookies)
	}
}

func TestJarGetCookiesEmptyHost(t *testing.T) {
	jar := NewJar()
	testURL, _ := url.Parse("https://example.com")

	jar.SetCookies(testURL, []*http.Cookie{})

	retrievedCookies := jar.Cookies(testURL)
	if len(retrievedCookies) != 0 {
		t.Errorf("Expected empty cookie slice, got %d cookies", len(retrievedCookies))
	}
}

func TestJarConcurrentAccess(t *testing.T) {
	jar := NewJar()
	testURL1, _ := url.Parse("https://example1.com")
	testURL2, _ := url.Parse("https://example2.com")

	cookies1 := []*http.Cookie{{Name: "test1", Value: "value1"}}
	cookies2 := []*http.Cookie{{Name: "test2", Value: "value2"}}

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			jar.SetCookies(testURL1, cookies1)
			jar.Cookies(testURL1)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			jar.SetCookies(testURL2, cookies2)
			jar.Cookies(testURL2)
		}
		done <- true
	}()

	<-done
	<-done

	stored1 := jar.Cookies(testURL1)
	stored2 := jar.Cookies(testURL2)

	if !reflect.DeepEqual(stored1, cookies1) {
		t.Error("Concurrent access affected cookies1")
	}

	if !reflect.DeepEqual(stored2, cookies2) {
		t.Error("Concurrent access affected cookies2")
	}
}

func TestJarHttpCookieIntegration(t *testing.T) {
	jar := NewJar()
	testURL, _ := url.Parse("https://example.com/path")

	httpCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "secret123",
		Path:     "/",
		Domain:   "example.com",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	jar.SetCookies(testURL, []*http.Cookie{httpCookie})

	retrievedCookies := jar.Cookies(testURL)
	if len(retrievedCookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(retrievedCookies))
	}

	retrieved := retrievedCookies[0]
	if retrieved.Name != httpCookie.Name {
		t.Errorf("Expected cookie name %s, got %s", httpCookie.Name, retrieved.Name)
	}
	if retrieved.Value != httpCookie.Value {
		t.Errorf("Expected cookie value %s, got %s", httpCookie.Value, retrieved.Value)
	}
	if retrieved.Path != httpCookie.Path {
		t.Errorf("Expected cookie path %s, got %s", httpCookie.Path, retrieved.Path)
	}
	if retrieved.Domain != httpCookie.Domain {
		t.Errorf("Expected cookie domain %s, got %s", httpCookie.Domain, retrieved.Domain)
	}
	if retrieved.MaxAge != httpCookie.MaxAge {
		t.Errorf("Expected cookie MaxAge %d, got %d", httpCookie.MaxAge, retrieved.MaxAge)
	}
	if retrieved.Secure != httpCookie.Secure {
		t.Errorf("Expected cookie Secure %t, got %t", httpCookie.Secure, retrieved.Secure)
	}
	if retrieved.HttpOnly != httpCookie.HttpOnly {
		t.Errorf("Expected cookie HttpOnly %t, got %t", httpCookie.HttpOnly, retrieved.HttpOnly)
	}
	if retrieved.SameSite != httpCookie.SameSite {
		t.Errorf("Expected cookie SameSite %v, got %v", httpCookie.SameSite, retrieved.SameSite)
	}
}
