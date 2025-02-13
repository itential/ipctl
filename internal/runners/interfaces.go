// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

type Runner interface {
}

type RunnerFunc func(Request) (*Response, error)

type Copier interface {
	Copy(Request) (*Response, error)
	CopyFrom(string, string) (any, error)
	CopyTo(string, any, bool) (any, error)
}

type Reader interface {
	Get(Request) (*Response, error)
	Describe(Request) (*Response, error)
}

type Writer interface {
	Create(Request) (*Response, error)
	Delete(Request) (*Response, error)
	Clear(Request) (*Response, error)
}

type Editor interface {
	Edit(Request) (*Response, error)
}

type Importer interface {
	Import(Request) (*Response, error)
}

type Exporter interface {
	Export(Request) (*Response, error)
}

type Controller interface {
	Start(Request) (*Response, error)
	Stop(Request) (*Response, error)
	Restart(Request) (*Response, error)
}

type Shower interface {
	Show(Request) (*Response, error)
}

type Inspector interface {
	Inspect(Request) (*Response, error)
}

type Gitter interface {
	Pull(Request) (*Response, error)
	Push(Request) (*Response, error)
}
