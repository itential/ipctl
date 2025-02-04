// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"embed"
	"strings"

	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/pkg/logger"
)

const descriptorsDir = "descriptors"

const (
	apiDescriptor = "api"

	accountsDescriptor  = "accounts"
	groupsDescriptor    = "groups"
	modelsDescriptor    = "models"
	rolesDescriptor     = "roles"
	roleTypesDescriptor = "roletypes"
	adaptersDescriptor  = "adapters"
	methodsDescriptor   = "methods"
	viewsDescriptor     = "views"
	bundleDescriptor    = "bundle"
	prebuiltsDescriptor = "prebuilts"
	profilesDescriptor  = "profiles"
	tagsDescriptor      = "tags"
	integrationModels   = "integration_models"
	integrations        = "integrations"
	adapterModels       = "adapter_models"

	automationsDescriptor = "automations"

	commandTemplatesDescriptor  = "command_templates"
	workflowsDescriptor         = "workflows"
	transformationsDescriptor   = "transformations"
	jsonformsDescriptor         = "jsonforms"
	projectsDescriptor          = "projects"
	analyticTemplatesDescriptor = "analytic_templates"
	templatesDescriptor         = "templates"

	devicesDescriptor              = "devices"
	deviceGroupsDescriptor         = "devicegroups"
	configurationParsersDescriptor = "configuration_parsers"
	gctreesDescriptor              = "gctrees"

	serverDescriptor = "server"

	localAAADescriptor    = "localaaa"
	localClientDescriptor = "localclient"
)

//go:embed descriptors/*.yaml
var content embed.FS

type DescriptorMap map[string]cmdutils.Descriptor
type Descriptors map[string]DescriptorMap

func loadDescriptors() Descriptors {
	descriptors := map[string]DescriptorMap{}

	entries, err := content.ReadDir(descriptorsDir)
	if err != nil {
		logger.Fatal(err, "failed to read descriptors directory")
	}

	for _, ele := range entries {
		name := strings.Split(ele.Name(), ".")[0]
		fn := strings.Join([]string{descriptorsDir, ele.Name()}, "/")

		data, err := content.ReadFile(fn)
		if err != nil {
			logger.Fatal(err, "failed to read descriptor")
		}

		descriptors[name] = cmdutils.LoadDescriptor(data)
	}

	return descriptors
}
