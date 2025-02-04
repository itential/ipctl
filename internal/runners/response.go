// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"errors"

	"github.com/itential/ipctl/pkg/logger"
)

type ResponseOption func(r *Response)

type Response struct {
	reference any

	Text  string
	Lines []string
	Json  []byte
	Url   string
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

func WithUrl(s string) ResponseOption {
	return func(r *Response) {
		r.Url = s
	}
}

func WithJson(o any) ResponseOption {
	return func(r *Response) {
		b, err := json.MarshalIndent(o, "", "    ")
		if err != nil {
			logger.Fatal(err, "failed to marshal json")
		}
		r.Json = b
		r.reference = o
	}
}

func NotImplemented(in Request) (*Response, error) {
	return nil, errors.New("command not implemented!!")
}
