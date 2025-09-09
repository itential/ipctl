// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

func setupDeviceService() *DeviceService {
	return NewDeviceService(
		testlib.Setup(),
	)
}

func TestDeviceService_GetAll(t *testing.T) {
	svc := setupDeviceService()
	defer testlib.Teardown()

	// Mock response for devices
	mockResponse := `{"entity":"devices","total":1,"totalByAdapter":{"test-adapter":1},"unique_device_count":1,"return_count":1,"start_index":0,"list":[{"name":"test-device","host":"192.168.1.1","ostype":"ios","device-type":"router","actions":["backup","restore"],"origins":null,"custom_property":"value"}]}`

	testlib.AddPostResponseToMux("/configuration_manager/devices", mockResponse, 200)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	if len(res) > 0 {
		assert.Equal(t, "test-device", res[0].Name)
		assert.Equal(t, "192.168.1.1", res[0].Host)
		assert.Equal(t, "ios", res[0].OsType)
		assert.Equal(t, "router", res[0].DeviceType)
		assert.Contains(t, res[0].Actions, "backup")
		assert.NotNil(t, res[0].Properties)
		assert.Equal(t, "value", res[0].Properties["custom_property"])
	}
}

func TestDeviceService_Get(t *testing.T) {
	svc := setupDeviceService()
	defer testlib.Teardown()

	// Mock response for single device
	mockResponse := `{"entity":"devices","total":1,"totalByAdapter":{"test-adapter":1},"unique_device_count":1,"return_count":1,"start_index":0,"list":[{"name":"test-device","host":"192.168.1.1","ostype":"ios","device-type":"router","actions":["backup","restore"],"origins":null,"custom_property":"value"}]}`

	testlib.AddPostResponseToMux("/configuration_manager/devices", mockResponse, 200)

	res, err := svc.Get("test-device")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test-device", res.Name)
	assert.Equal(t, "192.168.1.1", res.Host)
	assert.Equal(t, "ios", res.OsType)
	assert.Equal(t, "router", res.DeviceType)
}

func TestDeviceService_Get_NotFound(t *testing.T) {
	svc := setupDeviceService()
	defer testlib.Teardown()

	// Mock response for no devices found
	mockResponse := `{"entity":"devices","total":0,"totalByAdapter":{},"unique_device_count":0,"return_count":0,"start_index":0,"list":[]}`

	testlib.AddPostResponseToMux("/configuration_manager/devices", mockResponse, 200)

	res, err := svc.Get("nonexistent-device")

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "device not found")
}

func TestDeviceService_unmarshal(t *testing.T) {
	svc := setupDeviceService()

	deviceMap := map[string]interface{}{
		"name":            "test-device",
		"host":            "192.168.1.1",
		"ostype":          "ios",
		"device-type":     "router",
		"actions":         []string{"backup", "restore"},
		"origins":         nil,
		"custom_property": "value",
		"another_field":   123,
	}

	var device Device
	err := svc.unmarshal(deviceMap, &device)

	assert.Nil(t, err)
	assert.Equal(t, "test-device", device.Name)
	assert.Equal(t, "192.168.1.1", device.Host)
	assert.Equal(t, "ios", device.OsType)
	assert.Equal(t, "router", device.DeviceType)
	assert.NotNil(t, device.Properties)
	assert.Equal(t, "value", device.Properties["custom_property"])
	assert.Equal(t, 123, device.Properties["another_field"])
}
