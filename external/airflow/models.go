package airflow

import "time"

type DagRunListResponse struct {
	DagRuns      []DagRun `json:"dag_runs"`
	TotalEntries int      `json:"total_entries"`
}

type DagRun struct {
	ExecutionDate          time.Time `json:"execution_date"`
	State                  string    `json:"state"`
	ExternalTrigger        bool      `json:"external_trigger"`
	DagRunID               string    `json:"dag_run_id"`
	DagID                  string    `json:"dag_id"`
	LogicalDate            time.Time `json:"logical_date"`
	StartDate              time.Time `json:"start_date"`
	EndDate                time.Time `json:"end_date"`
	DataIntervalStart      time.Time `json:"data_interval_start"`
	DataIntervalEnd        time.Time `json:"data_interval_end"`
	LastSchedulingDecision time.Time `json:"last_scheduling_decision"`
	RunType                string    `json:"run_type"`
}

type DagRunRequest struct {
	OrderBy          string   `json:"order_by"`
	PageOffset       int      `json:"page_offset"`
	PageLimit        int      `json:"page_limit"`
	DagIds           []string `json:"dag_ids"` // nolint: revive
	ExecutionDateGte string   `json:"execution_date_gte,omitempty"`
	ExecutionDateLte string   `json:"execution_date_lte,omitempty"`
}

type DAGs struct {
	DAGS         []DAGObj `json:"dags"`
	TotalEntries int      `json:"total_entries"`
}

type DAGObj struct {
	DAGDisplayName              string   `json:"dag_display_name"`
	DAGID                       string   `json:"dag_id"`
	DefaultView                 string   `json:"default_view"`
	Description                 *string  `json:"description"`
	FileToken                   string   `json:"file_token"`
	Fileloc                     string   `json:"fileloc"`
	HasImportErrors             bool     `json:"has_import_errors"`
	HasTaskConcurrencyLimits    bool     `json:"has_task_concurrency_limits"`
	IsActive                    bool     `json:"is_active"`
	IsPaused                    bool     `json:"is_paused"`
	IsSubdag                    bool     `json:"is_subdag"`
	LastExpired                 *string  `json:"last_expired"`
	LastParsedTime              string   `json:"last_parsed_time"`
	LastPickled                 *string  `json:"last_pickled"`
	MaxActiveRuns               int      `json:"max_active_runs"`
	MaxActiveTasks              int      `json:"max_active_tasks"`
	MaxConsecutiveFailedDAGRuns int      `json:"max_consecutive_failed_dag_runs"`
	NextDagRun                  string   `json:"next_dagrun"`
	NextDagRunCreateAfter       string   `json:"next_dagrun_create_after"`
	NextDagRunDataIntervalEnd   string   `json:"next_dagrun_data_interval_end"`
	NextDagRunDataIntervalStart string   `json:"next_dagrun_data_interval_start"`
	Owners                      []string `json:"owners"`
	PickleID                    *string  `json:"pickle_id"`
	RootDagID                   *string  `json:"root_dag_id"`
	ScheduleInterval            Schedule `json:"schedule_interval"`
	SchedulerLock               *string  `json:"scheduler_lock"`
	Tags                        []Tag    `json:"tags"`
	TimetableDescription        string   `json:"timetable_description"`
}

type Schedule struct {
	Type  string `json:"__type"`
	Value string `json:"value"`
}

type Tag struct {
	Name string `json:"name"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	TaskId                  string    `json:"task_id"`
	TaskDisplayName         string    `json:"task_display_name"`
	Owner                   string    `json:"owner"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	TriggerRule             string    `json:"trigger_rule"`
	DependsOnPast           bool      `json:"depends_on_past"`
	IsMapped                bool      `json:"is_mapped"`
	WaitForDownstream       bool      `json:"wait_for_downstream"`
	Retries                 int       `json:"retries"`
	Queue                   string    `json:"queue"`
	Executor                string    `json:"executor"`
	Pool                    string    `json:"pool"`
	PoolSlots               int       `json:"pool_slots"`
	RetryExponentialBackoff bool      `json:"retry_exponential_backoff"`
	PriorityWeight          int       `json:"priority_weight"`
	WeightRule              string    `json:"weight_rule"`
	DownstreamTaskIds       []string  `json:"downstream_task_ids"`
	DocMd                   string    `json:"doc_md"`
}

type TaskInstances struct {
	TaskInstances []TaskInstance `json:"task_instances"`
	TotalEntries  int            `json:"total_entries"`
}

type TaskInstance struct {
	TaskId          string    `json:"task_id"`
	TaskDisplayName string    `json:"task_display_name"`
	DagId           string    `json:"dag_id"`
	DagRunId        string    `json:"dag_run_id"`
	ExecutionDate   time.Time `json:"execution_date"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Duration        float64   `json:"duration"`
	State           string    `json:"state"`
	TryNumber       int       `json:"try_number"`
	Note            string    `json:"note"`
}
