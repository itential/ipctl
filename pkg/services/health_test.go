// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	healthStatusSuccess      = "health/status.success.json"
	healthSystemSuccess      = "health/system.success.json"
	healthServerSuccess      = "health/server.success.json"
	healthApplicationsSuccess = "health/applications.success.json"
	healthAdaptersSuccess    = "health/adapters.success.json"
)

func setupHealthService() *HealthService {
	return NewHealthService(
		testlib.Setup(),
	)
}

func TestNewHealthService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewHealthService(client)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
	assert.Equal(t, reflect.TypeOf((*HealthService)(nil)), reflect.TypeOf(svc))
}

func TestHealthServiceGetStatus(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, healthStatusSuccess),
		)
		testlib.AddGetResponseToMux("/health/status", response, 0)

		res, err := svc.GetStatus()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*HealthStatus)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "test-server.example.com", res.Host)
		assert.Equal(t, "server-123", res.ServerId)
		assert.Equal(t, "Test Itential Platform", res.ServerName)
		assert.Equal(t, 2, len(res.Services))
		assert.Equal(t, "automation-studio", res.Services[0].Service)
		assert.Equal(t, "running", res.Services[0].Status)
		assert.Equal(t, "workflow-builder", res.Services[1].Service)
		assert.Equal(t, "running", res.Services[1].Status)
		assert.Equal(t, 1640995200, res.Timestamp)
		assert.Equal(t, "running", res.Apps)
		assert.Equal(t, "running", res.Adapters)
	}
}

func TestHealthServiceGetStatusError(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/health/status", "", 0)

	res, err := svc.GetStatus()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestHealthServiceGetSystemHealth(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, healthSystemSuccess),
		)
		testlib.AddGetResponseToMux("/health/system", response, 0)

		res, err := svc.GetSystemHealth()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*SystemHealth)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "x64", res.Arch)
		assert.Equal(t, "Linux 5.4.0-74-generic", res.Release)
		assert.Equal(t, 86400.5, res.Uptime)
		assert.Equal(t, int64(4294967296), res.FreeMem)
		assert.Equal(t, int64(17179869184), res.TotalMem)
		assert.Equal(t, 3, len(res.LoadAvg))
		assert.Equal(t, float32(0.1), res.LoadAvg[0])
		assert.Equal(t, 2, len(res.Cpus))
		assert.Equal(t, "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz", res.Cpus[0].Model)
		assert.Equal(t, 2600, res.Cpus[0].Speed)
		assert.Equal(t, 1000, res.Cpus[0].Times.User)
		assert.Equal(t, 500, res.Cpus[0].Times.Sys)
		assert.Equal(t, 98500, res.Cpus[0].Times.Idle)
	}
}

func TestHealthServiceGetSystemHealthError(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/health/system", "", 0)

	res, err := svc.GetSystemHealth()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestHealthServiceGetServerHealth(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, healthServerSuccess),
		)
		testlib.AddGetResponseToMux("/health/server", response, 0)

		res, err := svc.GetServerHealth()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*ServerHealth)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "18.16.0", res.Version)
		assert.Equal(t, "node", res.Release)
		assert.Equal(t, "x64", res.Arch)
		assert.Equal(t, "linux", res.Platform)
		assert.NotEmpty(t, res.Versions)
		assert.Equal(t, "18.16.0", res.Versions["node"])
		assert.Equal(t, int32(67108864), res.MemoryUsage.Rss)
		assert.Equal(t, int32(33554432), res.MemoryUsage.HeapTotal)
		assert.Equal(t, int32(25165824), res.MemoryUsage.HeapUsed)
		assert.Equal(t, int64(150000), res.CpuUsage.User)
		assert.Equal(t, int64(50000), res.CpuUsage.System)
		assert.Equal(t, 3600.25, res.Uptime)
		assert.Equal(t, 12345, res.Pid)
		assert.NotEmpty(t, res.Dependencies)
		assert.Equal(t, "2023.2.9", res.Dependencies["itential-platform"])
	}
}

func TestHealthServiceGetServerHealthError(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/health/server", "", 0)

	res, err := svc.GetServerHealth()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestHealthServiceGetApplicationHealth(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, healthApplicationsSuccess),
		)
		testlib.AddGetResponseToMux("/health/applications", response, 0)

		res, err := svc.GetApplicationHealth()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
		
		// Test first application
		app1 := res[0]
		assert.Equal(t, "automation-studio", app1.Id)
		assert.Equal(t, "@itential/automation-studio", app1.PackageId)
		assert.Equal(t, "4.15.2", app1.Version)
		assert.Equal(t, "application", app1.Type)
		assert.Equal(t, "Itential Automation Studio", app1.Description)
		assert.Equal(t, "/automation-studio", app1.RoutePrefix)
		assert.Equal(t, "running", app1.State)
		assert.Equal(t, "connected", app1.Connection.Status)
		assert.Equal(t, 3600.123, app1.Uptime)
		assert.Equal(t, int32(134217728), app1.MemoryUsage.Rss)
		assert.Equal(t, int64(250000), app1.CpuUsage.User)
		assert.Equal(t, float64(12346), app1.Pid)
		assert.Equal(t, "error", app1.Logger.Console)
		assert.Equal(t, "info", app1.Logger.File)
		assert.Equal(t, int64(1640995200000), app1.Timestamp)

		// Test second application
		app2 := res[1]
		assert.Equal(t, "workflow-builder", app2.Id)
		assert.Equal(t, "@itential/workflow-builder", app2.PackageId)
		assert.Equal(t, "2.48.0", app2.Version)
		assert.Equal(t, "application", app2.Type)
		assert.Equal(t, "Itential Workflow Builder", app2.Description)
		assert.Equal(t, "/workflow-builder", app2.RoutePrefix)
		assert.Equal(t, "running", app2.State)
		assert.Equal(t, "connected", app2.Connection.Status)
	}
}

func TestHealthServiceGetApplicationHealthError(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/health/applications", "", 0)

	res, err := svc.GetApplicationHealth()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestHealthServiceGetAdapterHealth(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, healthAdaptersSuccess),
		)
		testlib.AddGetResponseToMux("/health/adapters", response, 0)

		res, err := svc.GetAdapterHealth()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
		
		// Test first adapter
		adapter1 := res[0]
		assert.Equal(t, "adapter-netmiko", adapter1.Id)
		assert.Equal(t, "@itential/adapter-netmiko", adapter1.PackageId)
		assert.Equal(t, "1.2.3", adapter1.Version)
		assert.Equal(t, "adapter", adapter1.Type)
		assert.Equal(t, "Netmiko Adapter", adapter1.Description)
		assert.Equal(t, "/adapter-netmiko", adapter1.RoutePrefix)
		assert.Equal(t, "running", adapter1.State)
		assert.Equal(t, "connected", adapter1.Connection.Status)
		assert.Equal(t, 3550.789, adapter1.Uptime)
		assert.Equal(t, int32(83886080), adapter1.MemoryUsage.Rss)
		assert.Equal(t, int64(120000), adapter1.CpuUsage.User)
		assert.Equal(t, float64(12348), adapter1.Pid)
		assert.Equal(t, "info", adapter1.Logger.Console)
		assert.Equal(t, "debug", adapter1.Logger.File)

		// Test second adapter
		adapter2 := res[1]
		assert.Equal(t, "adapter-ansible", adapter2.Id)
		assert.Equal(t, "@itential/adapter-ansible", adapter2.PackageId)
		assert.Equal(t, "2.1.0", adapter2.Version)
		assert.Equal(t, "adapter", adapter2.Type)
		assert.Equal(t, "Ansible Adapter", adapter2.Description)
		assert.Equal(t, "/adapter-ansible", adapter2.RoutePrefix)
		assert.Equal(t, "running", adapter2.State)
		assert.Equal(t, "connected", adapter2.Connection.Status)
	}
}

func TestHealthServiceGetAdapterHealthError(t *testing.T) {
	svc := setupHealthService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/health/adapters", "", 0)

	res, err := svc.GetAdapterHealth()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}