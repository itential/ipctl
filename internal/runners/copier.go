// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

type CopyRequest struct {
	Request Request
	Type    string
}

type CopyResponse struct {
	Name string
	From string
	To   string
}

func Copy(in CopyRequest, r Copier) (*CopyResponse, error) {
	logger.Trace()

	name := in.Request.Args[0]

	var common *flags.AssetCopyCommon
	utils.LoadObject(in.Request.Common, &common)

	if common.From == common.To {
		return nil, errors.New("source (--from) and destination (--to) servers must be different values")
	}

	logger.Info("attempting to copy `%s` (type of %s) from `%s` to `%s`", name, in.Type, common.From, common.To)

	src, err := r.CopyFrom(common.From, name)
	if err != nil {
		return nil, err
	}

	_, err = r.CopyTo(common.To, src, common.Replace)

	return &CopyResponse{
		Name: name,
		From: common.From,
		To:   common.To,
	}, nil
}
