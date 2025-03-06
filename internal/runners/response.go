// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/pkg/config"
)

type ResponseOption func(r *Response)

type Response struct {
	// Object is a generate interface field that holds the response object.
	Object any

	Text  string
	Lines []string
}

func (r *Response) String() string {
	return r.Text
}

func NewResponse(text string, opts ...ResponseOption) *Response {
	r := &Response{Text: text}
	for _, ele := range opts {
		ele(r)
	}
	return r
}

func WithTable(lines []string) ResponseOption {
	return func(r *Response) {
		r.Lines = lines
	}
}

func WithObject(o any) ResponseOption {
	return func(r *Response) {
		r.Object = o
	}
}

func NotImplemented(in Request) (*Response, error) {
	return nil, errors.New("command not implemented!!")
}

func makeUrl(cfg *config.Config, path string, args ...any) string {
	profile, _ := cfg.ActiveProfile()

	var u string

	if profile.UseTLS {
		u = "https://"
	} else {
		u = "http://"
	}

	u += profile.Host

	if profile.Port != 0 {
		u += fmt.Sprintf(":%v", profile.Port)
	}

	path = fmt.Sprintf(path, args...)

	if strings.HasPrefix(path, "/") {
		u += path
	} else {
		u += fmt.Sprintf("/%s", path)
	}

	return u
}
