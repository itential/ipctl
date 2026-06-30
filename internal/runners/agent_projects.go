// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/resources"
	"github.com/itential/ipctl/pkg/services"
)

// AgentProjectRunner orchestrates CLI commands for agent project management.
// It implements the Reader, Importer, and Exporter interfaces.
type AgentProjectRunner struct {
	BaseRunner
	resource resources.AgentProjectResourcer
}

// NewAgentProjectRunner creates a new AgentProjectRunner with the provided client and config.
func NewAgentProjectRunner(c client.Client, cfg config.Provider) *AgentProjectRunner {
	return &AgentProjectRunner{
		BaseRunner: NewBaseRunner(c, cfg),
		resource:   resources.NewAgentProjectResource(services.NewAgentProjectService(c)),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get retrieves all agent projects and returns them for display.
func (r *AgentProjectRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	projects, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: projects,
	}, nil
}

// Describe retrieves detailed information about a specific agent project.
func (r *AgentProjectRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	project, err := r.resource.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	createdBy := extractAgentProjectUsername(project.CreatedBy, "unknown")
	updatedBy := extractAgentProjectUsername(project.LastUpdatedBy, "unknown")

	output := []string{
		fmt.Sprintf("Name: %s (%s)", project.Name, project.Id),
		fmt.Sprintf("Description: %s", project.Description),
		fmt.Sprintf("Created: %s, by: %s", project.Created, createdBy),
		fmt.Sprintf("Updated: %s, by: %s", project.LastUpdated, updatedBy),
		fmt.Sprintf("Components: %d", len(project.Components)),
	}

	return &Response{
		Text:   strings.Join(output, "\n"),
		Object: project,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import imports an agent project bundle from a local file or Git repository.
func (r *AgentProjectRunner) Import(in Request) (*Response, error) {
	logging.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return nil, err
	}

	wd := filepath.Dir(path)

	if common.Repository != "" {
		defer os.RemoveAll(wd)
	}

	var bundle services.AgentProjectBundle

	if err := importLoadFromDisk(path, &bundle); err != nil {
		return nil, err
	}

	if !common.Replace {
		existing, err := r.resource.GetByName(bundle.Name)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("agent project %q already exists, use --replace to overwrite", bundle.Name)
		}
	}

	imported, err := r.resource.Import(bundle)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported agent project `%s` (%s)", imported.Name, imported.Id),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export exports an agent project bundle to a local file or Git repository.
func (r *AgentProjectRunner) Export(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	project, err := r.resource.GetByName(name)
	if err != nil {
		return nil, err
	}

	bundle, err := r.resource.Export(project.Id)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(bundle)
	if err != nil {
		return nil, err
	}

	var exported map[string]interface{}
	if err := json.Unmarshal(b, &exported); err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.agent-project.json", normalizeFilename(name))

	if err := exportAssetFromRequest(in, exported, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported agent project `%s`", bundle.Name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private helpers
//

// extractAgentProjectUsername safely extracts a username from a user object.
func extractAgentProjectUsername(userObj any, fallback string) string {
	if userObj == nil {
		return fallback
	}

	userMap, ok := userObj.(map[string]interface{})
	if !ok {
		return fallback
	}

	username, ok := userMap["username"].(string)
	if !ok {
		return fallback
	}

	return username
}
