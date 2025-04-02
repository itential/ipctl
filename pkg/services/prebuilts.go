// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type PrebuiltPackageMetadataRepository struct {
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
}

type PrebuiltPackageMetadata struct {
	Name         string                            `json:"name"`
	Version      string                            `json:"version"`
	Description  string                            `json:"description"`
	License      string                            `json:"license"`
	Repository   PrebuiltPackageMetadataRepository `json:"repository"`
	Keywords     []string                          `json:"keywords"`
	Author       string                            `json:"author"`
	Dependencies map[string]interface{}            `json:"IAPDependencies"`
	GitlabId     int                               `json:"gitlabId"`
}

type PrebuiltPackageManifest struct {
	Name        string                            `json:"bundleName"`
	Fingerprint string                            `json:"fingerprint"`
	Epoch       string                            `json:"createdEpoch"`
	Artifacts   []PrebuiltPackageManifestArtifact `json:"artifacts"`
}

type PrebuiltPackageManifestArtifact struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Location   string                 `json:"location"`
	Properties map[string]interface{} `json:"properties"`
}

type PrebuiltPackageBundle struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type PrebuiltPackage struct {
	Metadata PrebuiltPackageMetadata `json:"metadata"`
	Manifest PrebuiltPackageManifest `json:"manifest"`
	Bundles  []PrebuiltPackageBundle `json:"bundles"`
	Readme   string                  `json:"readme"`
}

type PrebuiltComponent struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type PrebuiltRepository struct {
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
	Type     string `json:"type"`
}

type PrebuiltManifestArtifact struct {
	Id         string                 `json:"id"`
	Location   string                 `json:"location"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
	Type       string                 `json:"type"`
}

type PrebuiltManifest struct {
	Artifacts    []PrebuiltManifestArtifact `json:"artifacts"`
	BundleName   string                     `json:"bundleName"`
	CreatedEpoch string                     `json:"createdEpoch"`
	Fingerprint  string                     `json:"fingerprint"`
}

type Prebuilt struct {
	Id           string                 `json:"_id"`
	Name         string                 `json:"name"`
	Dependencies map[string]interface{} `json:"IAPDependencies"`
	Author       string                 `json:"author"`
	Components   []PrebuiltComponent    `json:"components"`
	Description  string                 `json:"description"`
	Installed    string                 `json:"installed"`
	InstalledBy  string                 `json:"installedBy"`
	Keywords     []string               `json:"keywords"`
	License      string                 `json:"license"`
	Manifest     PrebuiltManifest       `json:"manifest"`
	Readme       string                 `json:"readme"`
	Repository   PrebuiltRepository     `json:"repository"`
	Version      string                 `json:"version"`
}

type PrebuiltService struct {
	client client.Client
}

func NewPrebuiltService(c client.Client) *PrebuiltService {
	return &PrebuiltService{client: c}
}

func (svc *PrebuiltService) Get(id string) (*Prebuilt, error) {
	logger.Trace()

	res, err := Do(&Request{
		client: svc.client,
		method: http.MethodGet,
		uri:    fmt.Sprintf("/prebuilts/%s", id),
	})

	if err != nil {
		return nil, err
	}

	var prebuilt *Prebuilt

	if err := json.Unmarshal(res.Body, &prebuilt); err != nil {
		return nil, err
	}

	return prebuilt, nil
}

func (p *PrebuiltService) GetByName(name string) (*Prebuilt, error) {
	logger.Trace()

	prebuilts, err := p.GetAll()
	if err != nil {
		return nil, err
	}

	var prebuiltId string
	for _, ele := range prebuilts {
		if ele.Name == name {
			prebuiltId = ele.Id
		}
	}

	if prebuiltId == "" {
		return nil, errors.New(fmt.Sprintf("prebuilt `%s` does not exist", name))
	}

	prebuilt, err := p.Get(prebuiltId)
	if err != nil {
		return nil, err
	}

	return prebuilt, nil
}

func (svc *PrebuiltService) GetAll() ([]Prebuilt, error) {
	logger.Trace()

	res, err := Do(&Request{
		client: svc.client,
		method: http.MethodGet,
		uri:    "/prebuilts",
	})

	if err != nil {
		return nil, err
	}

	type Response struct {
		Results []Prebuilt `json:"results"`
		Total   int        `json:"total"`
	}

	var response Response

	if err := json.Unmarshal(res.Body, &response); err != nil {
		return nil, err
	}

	logger.Info("Found %v prebuilts", response.Total)

	return response.Results, nil
}

func (svc *PrebuiltService) Delete(id string) error {
	logger.Trace()

	_, err := Do(&Request{
		client: svc.client,
		method: http.MethodDelete,
		uri:    fmt.Sprintf("/prebuilts/%s", id),
	})

	return err
}

func (svc *PrebuiltService) Import(in PrebuiltPackage, overwrite bool) (*Prebuilt, error) {
	logger.Trace()

	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	var prebuilt map[string]interface{}
	if err := json.Unmarshal(b, &prebuilt); err != nil {
		return nil, err
	}

	return svc.ImportRaw(prebuilt, overwrite)
}

func (svc *PrebuiltService) ImportRaw(in any, overwrite bool) (*Prebuilt, error) {
	logger.Trace()

	body := map[string]interface{}{
		"options":  map[string]interface{}{"overwrite": overwrite},
		"prebuilt": in,
	}

	res, err := Do(&Request{
		client:             svc.client,
		method:             http.MethodPost,
		uri:                "/prebuilts/import",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	})
	if err != nil {
		return nil, err
	}

	type Response struct {
		Message string    `json:"message"`
		Status  string    `json:"status"`
		Data    *Prebuilt `json:"data"`
	}

	var response Response

	if err := json.Unmarshal(res.Body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (svc *PrebuiltService) Export(id string) (*PrebuiltPackage, error) {
	logger.Trace()

	res, err := Do(&Request{
		client: svc.client,
		method: http.MethodGet,
		uri:    fmt.Sprintf("/prebuilts/%s/export", id),
	})

	if err != nil {
		return nil, err
	}

	var pb *PrebuiltPackage
	if err := json.Unmarshal(res.Body, &pb); err != nil {
		return nil, err
	}

	return pb, nil
}
