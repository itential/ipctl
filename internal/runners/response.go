// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/config"
)

type ResponseOption func(r *Response)

// Response is a return object from a runner that provides the necessary
// elements to the Handler for displaying information.
type Response struct {
	Object   any
	Text     string
	Template string
	Keys     []string
}

// String implements the Stringer interface.  This function will return the
// Response object as a string.
func (r *Response) String() string {
	if len(r.Keys) > 0 {
		table, err := renderTable(r.Object, r.Keys)
		if err != nil {
			logging.Fatal(err, "")
		}
		return strings.Join(table, "\n")
	} else if r.Template != "" {
		tmpl, err := template.New("output").Parse(r.Template)
		if err != nil {
			logging.Fatal(err, "")
		}

		var b bytes.Buffer
		if err := tmpl.Execute(&b, r.Object); err != nil {
			logging.Fatal(err, "")
		}

		return b.String()
	} else if r.Text == "" && r.Object != nil {
		b, err := json.MarshalIndent(r.Object, "", "    ")
		if err != nil {
			logging.Fatal(err, "")
		}
		return string(b)
	} else if r.Text != "" {
		return r.Text
	} else {
		return "error formating response object"
	}
}

func notImplemented(in Request) (*Response, error) {
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

func renderTable(in any, keys []string) ([]string, error) {
	data, err := toArrayOfMaps(in)
	if err != nil {
		return nil, err
	}

	var headers []string

	for _, ele := range keys {
		headers = append(headers, strings.ToUpper(ele))
	}

	var output = []string{strings.Join(headers, "\t")}

	for _, ele := range data {
		var row []string
		for _, k := range keys {
			if v, exists := ele[k]; exists {
				row = append(row, v.(string))
			} else {
				terminal.Warning("table field `%s` does not exist", k)
			}

		}
		output = append(output, strings.Join(row, "\t"))
	}

	return output, nil
}
