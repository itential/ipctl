// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

const (
	exportSuccessMessage = "Successfully exported %s `%s` to `%s`"
)

type ExportAction struct {
	Filename string
	Data     any
	Common   *flags.AssetExportCommon
}

type ExportActionResponse struct {
	Filename string
	Message  string
}

func NewExportAction(data any, fn string, common *flags.AssetExportCommon) ExportAction {
	logger.Trace()
	return ExportAction{
		Data:     data,
		Filename: fn,
		Common:   common,
	}
}

func (a ExportAction) Do() error {
	logger.Trace()

	var common flags.AssetExportCommon
	utils.LoadObject(a.Common, &common)

	if err := utils.WriteJsonToDisk(a.Data, a.Filename, common.Path); err != nil {
		return err
	}

	return nil
}
