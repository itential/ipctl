// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "github.com/spf13/cobra"

// AgentProjectImportOptions holds options for the `import agent-project` command.
type AgentProjectImportOptions struct{}

func (o *AgentProjectImportOptions) Flags(_ *cobra.Command) {}

// AgentProjectExportOptions holds options for the `export agent-project` command.
type AgentProjectExportOptions struct{}

func (o *AgentProjectExportOptions) Flags(_ *cobra.Command) {}
