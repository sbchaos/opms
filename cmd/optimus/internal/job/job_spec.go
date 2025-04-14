package job

import (
	"time"
)

const NewWindowVersion = 3

type YamlSpec struct {
	Version      int                 `yaml:"version,omitempty"`
	Name         string              `yaml:"name"`
	Owner        string              `yaml:"owner,omitempty"`
	Description  string              `yaml:"description,omitempty"`
	Schedule     JobSpecSchedule     `yaml:"schedule,omitempty"`
	Behavior     JobSpecBehavior     `yaml:"behavior,omitempty"`
	Task         JobSpecTask         `yaml:"task,omitempty"`
	Asset        map[string]string   `yaml:"-"`
	Labels       map[string]string   `yaml:"labels,omitempty"`
	Hooks        []JobSpecHook       `yaml:"hooks,omitempty"`
	Dependencies []JobSpecDependency `yaml:"dependencies,omitempty"`
	Metadata     *JobSpecMetadata    `yaml:"metadata,omitempty"`
	Path         string              `yaml:"-"`
}

type JobSpecSchedule struct {
	StartDate string `yaml:"start_date,omitempty"`
	EndDate   string `yaml:"end_date,omitempty"`
	Interval  string `yaml:"interval"`
}

type JobSpecBehavior struct {
	DependsOnPast bool                      `yaml:"depends_on_past"`
	Catchup       bool                      `yaml:"catch_up,omitempty"`
	Retry         *JobSpecBehaviorRetry     `yaml:"retry,omitempty"`
	Notify        []JobSpecBehaviorNotifier `yaml:"notify,omitempty"`
	Webhook       []JobSpecBehaviorWebhook  `yaml:"webhook,omitempty"`
}

type JobSpecBehaviorRetry struct {
	Count              int           `yaml:"count,omitempty"`
	Delay              time.Duration `yaml:"delay,omitempty"`
	ExponentialBackoff bool          `yaml:"exponential_backoff,omitempty"`
}

type JobSpecBehaviorNotifier struct {
	On       string            `yaml:"on,omitempty"`
	Config   map[string]string `yaml:"config,omitempty"`
	Channels []string          `yaml:"channels,omitempty"`
	Severity string            `yaml:"severity,omitempty"`
	Team     string            `yaml:"team,omitempty"`
}

type WebhookEndpoint struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type JobSpecBehaviorWebhook struct {
	On        string            `yaml:"on"`
	Endpoints []WebhookEndpoint `yaml:"endpoints"`
}

type JobSpecTask struct {
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config,omitempty"`
	Window JobSpecTaskWindow `yaml:"window,omitempty"`
}

type JobSpecTaskWindow struct {
	Size string `yaml:"size,omitempty"`
	// deprecated, replaced by ShiftBy
	Offset     string `yaml:"offset,omitempty"`
	TruncateTo string `yaml:"truncate_to,omitempty"`
	Preset     string `yaml:"preset,omitempty"`
	ShiftBy    string `yaml:"shift_by,omitempty"`
	Location   string `yaml:"location,omitempty"`
}

type JobSpecHook struct {
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config,omitempty"`
}

type JobSpecDependency struct {
	JobName string                 `yaml:"job,omitempty"`
	Type    string                 `yaml:"type,omitempty"`
	HTTP    *JobSpecDependencyHTTP `yaml:"http,omitempty"`
}

type JobSpecDependencyHTTP struct {
	Name          string            `yaml:"name"`
	RequestParams map[string]string `yaml:"params,omitempty"`
	URL           string            `yaml:"url"`
	Headers       map[string]string `yaml:"headers,omitempty"`
}

type JobSpecMetadata struct {
	Resource *JobSpecMetadataResource `yaml:"resource,omitempty"`
	Airflow  *JobSpecMetadataAirflow  `yaml:"airflow,omitempty"`
}

type JobSpecMetadataResource struct {
	Request *JobSpecMetadataResourceConfig `yaml:"request,omitempty"`
	Limit   *JobSpecMetadataResourceConfig `yaml:"limit,omitempty"`
}

type JobSpecMetadataResourceConfig struct {
	Memory string `yaml:"memory,omitempty"`
	CPU    string `yaml:"cpu,omitempty"`
}

type JobSpecMetadataAirflow struct {
	Pool  string `yaml:"pool" json:"pool,omitempty"`
	Queue string `yaml:"queue" json:"queue,omitempty"`
}
