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
