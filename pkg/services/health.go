// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ServiceStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type MemoryUsage struct {
	Rss          int32 `json:"rss"`
	HeapTotal    int32 `json:"heapTotal"`
	HeapUsed     int32 `json:"heapUsed"`
	External     int32 `json:"external"`
	ArrayBuffers int32 `json:"arrayBuffers"`
}

type CpuUsage struct {
	User   int64 `json:"user"`
	System int64 `json:"system"`
}

type ServerHealth struct {
	Version      string            `json:"version"`
	Release      string            `json:"release"`
	Arch         string            `json:"arch"`
	Platform     string            `json:"platform"`
	Versions     map[string]string `json:"versions"`
	MemoryUsage  MemoryUsage       `json:"memoryUsage"`
	CpuUsage     CpuUsage          `json:"cpuUsage"`
	Uptime       float64           `json:"uptime"`
	Pid          int               `json:"pid"`
	Dependencies map[string]string `json:"dependencies"`
}

type CpuTime struct {
	User int `json:"user"`
	Nice int `json:"nice"`
	Sys  int `json:"sys"`
	Idle int `json:"idle"`
	Irq  int `json:"irq"`
}

type Cpu struct {
	Model string  `json:"model"`
	Speed int     `json:"speed"`
	Times CpuTime `json:"times"`
}

type Connection struct {
	Status string `json:"state"`
}

type ApplicationHealth struct {
	Id          string      `json:"id"`
	PackageId   string      `json:"package_id"`
	Version     string      `json:"version"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	RoutePrefix string      `json:"routePrefix"`
	State       string      `json:"state"`
	Connection  Connection  `json:"connection"`
	Uptime      any         `json:"uptime"`
	MemoryUsage MemoryUsage `json:"memoryUsage"`
	CpuUsage    CpuUsage    `json:"cpuUsage"`
	Pid         any         `json:"pid"`
	Logger      Logger      `json:"logger"`
	Timestamp   int64       `json:"timestamp"`
	PrevUptime  float64     `json:"prevUptime"`
}

type Logger struct {
	Console string `json:"console"`
	File    string `json:"file"`
	Syslog  any    `json:"syslog"`
}

type SystemHealth struct {
	Arch     string    `json:"arch"`
	Release  string    `json:"release"`
	Uptime   float64   `json:"uptime"`
	FreeMem  int64     `json:"freemem"`
	TotalMem int64     `json:"totalmem"`
	LoadAvg  []float32 `json:"loadavg"`
	Cpus     []Cpu     `json:"cpus"`
}

type HealthStatus struct {
	Host       string          `json:"host"`
	ServerId   string          `json:"serverId"`
	ServerName string          `json:"serverName"`
	Services   []ServiceStatus `json:"services"`
	Timestamp  int             `json:"timestamp"`
	Apps       string          `json:"apps"`
	Adapters   string          `json:"adapters"`
}

type HealthService struct {
	BaseService
}

func NewHealthService(c client.Client) *HealthService {
	return &HealthService{BaseService: NewBaseService(c)}
}

func (svc *HealthService) GetStatus() (*HealthStatus, error) {
	logger.Trace()

	var res *HealthStatus

	if err := svc.BaseService.Get("/health/status", &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *HealthService) GetSystemHealth() (*SystemHealth, error) {
	logger.Trace()

	var res *SystemHealth

	if err := svc.BaseService.Get("/health/system", &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *HealthService) GetServerHealth() (*ServerHealth, error) {
	logger.Trace()

	var res *ServerHealth

	if err := svc.BaseService.Get("/health/server", &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *HealthService) GetApplicationHealth() ([]ApplicationHealth, error) {
	logger.Trace()

	type Response struct {
		Results []ApplicationHealth `json:"results"`
	}

	var res Response

	if err := svc.BaseService.Get("/health/applications", &res); err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (svc *HealthService) GetAdapterHealth() ([]ApplicationHealth, error) {
	logger.Trace()

	type Response struct {
		Results []ApplicationHealth `json:"results"`
	}

	var res Response

	if err := svc.BaseService.Get("/health/adapters", &res); err != nil {
		return nil, err
	}

	return res.Results, nil
}
