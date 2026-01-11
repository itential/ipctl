// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
)

type CopyRequest struct {
	Request Request
	Type    string
}

type CopyResponse struct {
	Name         string
	From         string
	CopyFromData any
	To           string
	CopyToData   any
}

func Copy(in CopyRequest, r Copier) (*CopyResponse, error) {
	logging.Trace()

	name := in.Request.Args[0]

	common := in.Request.Common.(*flags.AssetCopyCommon)

	if common.From == common.To {
		return nil, errors.New("source (--from) and destination (--to) servers must be different values")
	}

	logging.Info("attempting to copy `%s` (type of %s) from `%s` to `%s`", name, in.Type, common.From, common.To)

	src, err := r.CopyFrom(common.From, name)
	if err != nil {
		return nil, err
	}

	res, err := r.CopyTo(common.To, src, common.Replace)
	if err != nil {
		return nil, err
	}

	return &CopyResponse{
		Name:         name,
		From:         common.From,
		CopyFromData: src,
		To:           common.To,
		CopyToData:   res,
	}, nil
}
