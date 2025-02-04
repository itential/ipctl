// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type AutomationCreateOptions struct {
	Description string
	Replace     bool
}

func (o *AutomationCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Description of the automation")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the exist automationif it exists")
}

type AutomationImportOptions struct {
	DisableComponentCheck   bool
	DisableGroupReadCheck   bool
	DisableGroupWriteCheck  bool
	DisableGroupExistsCheck bool
}

func (o *AutomationImportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.DisableComponentCheck, "disable-component-check", o.DisableComponentCheck, "Disable checking if component is accessible prior to importing automation")
	cmd.Flags().BoolVar(&o.DisableGroupReadCheck, "disable-group-read-check", o.DisableGroupReadCheck, "Do not check if current user is member of read group")
	cmd.Flags().BoolVar(&o.DisableGroupWriteCheck, "disable-group-write-check", o.DisableGroupWriteCheck, "Do not check if current user is member of write group")
	cmd.Flags().BoolVar(&o.DisableGroupExistsCheck, "disable-group-exists-check", o.DisableGroupExistsCheck, "Disable checking if the read and write gorups exist before importing")
}
