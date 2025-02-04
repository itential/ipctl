// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

type Metric struct {
	JobsComplete     int    `json:"jobsComplete"`
	TotalRunTime     int    `json:"totalRunTime"`
	SlaTargetsMissed int    `json:"slaTargetsMissed"`
	TotalManualTime  int    `json:"totalManualTime"`
	StartDate        string `json:"startDate"`
}

type JobMetrics struct {
	Id                string                 `json:"_id"`
	Workflow          map[string]interface{} `json:"workflow"`
	Metrics           []Metric               `json:"metrics"`
	PreAutomationTime int                    `json:"preAutomationTime"`
}
