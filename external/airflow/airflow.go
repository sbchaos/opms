package airflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	dagStatusBatchURL = "api/v1/dags/~/dagRuns/list"
	dagURL            = "api/v1/dags"
	dagRunClearURL    = "api/v1/dags/%s/clearTaskInstances"
	dagRunCreateURL   = "api/v1/dags/%s/dagRuns"
	dagRunModifyURL   = "api/v1/dags/%s/dagRuns/%s"
	airflowDateFormat = "2006-01-02T15:04:05+00:00"

	taskInstances = "api/v1/dags/%s/dagRuns/%s/taskInstances"
)

type Airflow struct {
	client *Client
	auth   Auth
}

type JobRunsCriteria struct {
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Filter      []string
	OnlyLastRun bool
}

func (s *Airflow) FetchJobRunBatch(ctx context.Context, jobQuery *JobRunsCriteria) (*DagRunListResponse, error) {
	dagRunRequest := getDagRunRequest(jobQuery)
	reqBody, err := json.Marshal(dagRunRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal dag run request: %w", err)
	}

	req := Request{
		Path:   dagStatusBatchURL,
		Method: http.MethodPost,
		Body:   reqBody,
	}

	resp, err := s.client.Invoke(ctx, req, s.auth)
	if err != nil {
		return nil, fmt.Errorf("unable to invoke dag runs: %w", err)
	}

	var dagRunList DagRunListResponse
	if err := json.Unmarshal(resp, &dagRunList); err != nil {
		return nil, fmt.Errorf("json error on parsing airflow dag runs: %s, %w", string(resp), err)
	}

	return &dagRunList, nil
}

//
//func (s *Airflow) GetJobRunsWithDetails(ctx context.Context, criteria *JobRunsCriteria) ([]*scheduler.JobRunWithDetails, error) {
//	resp, err := s.FetchJobRunBatch(ctx, criteria)
//	if err != nil {
//		return nil, fmt.Errorf("failure while fetching airflow dag runs: %w", err)
//	}
//	var dagRunList DagRunListResponse
//	if err := json.Unmarshal(resp, &dagRunList); err != nil {
//		return nil, fmt.Errorf(fmt.Sprintf("json error on parsing airflow dag runs: %s", string(resp)), err)
//	}
//
//	return getJobRunsWithDetails(dagRunList)
//}
//
//func (s *Airflow) GetJobRuns(ctx context.Context, criteria *JobRunsCriteria) ([]*scheduler.JobRunStatus, error) {
//	resp, err := s.FetchJobRunBatch(ctx, criteria)
//	if err != nil {
//		return nil, fmt.Errorf("failure while fetching airflow dag runs: %w", err)
//	}
//
//	var dagRunList DagRunListResponse
//	if err := json.Unmarshal(resp, &dagRunList); err != nil {
//		return nil, fmt.Errorf(fmt.Sprintf("json error on parsing airflow dag runs: %s", string(resp)), err)
//	}
//
//	return getJobRuns(dagRunList)
//}

func (s *Airflow) fetchJobs(ctx context.Context, offset int) (*DAGs, error) {
	params := url.Values{}
	params.Add("limit", "100") // Default and max is 100
	params.Add("order_by", "dag_id")
	params.Add("offset", strconv.Itoa(offset))
	req := Request{
		Path:   dagURL,
		Method: http.MethodGet,
		Query:  params.Encode(),
	}

	resp, err := s.client.Invoke(ctx, req, s.auth)
	if err != nil {
		return nil, err
	}

	var dagsInfo DAGs
	err = json.Unmarshal(resp, &dagsInfo)
	if err != nil {
		return nil, err
	}
	return &dagsInfo, nil
}

func (s *Airflow) FetchAllJobs(ctx context.Context) (*DAGs, error) {
	var offset int
	var allDags DAGs
	for {
		fetchResp, err := s.fetchJobs(ctx, offset)
		if err != nil {
			return nil, err
		}

		allDags.DAGS = append(allDags.DAGS, fetchResp.DAGS...)
		if len(allDags.DAGS) < fetchResp.TotalEntries {
			offset = len(allDags.DAGS)
			continue
		}
		break
	}
	return &allDags, nil
}

func getDagRunRequest(criteria *JobRunsCriteria) DagRunRequest {
	if criteria.OnlyLastRun {
		return DagRunRequest{
			OrderBy:    "-execution_date",
			PageOffset: 0,
			PageLimit:  1,
			DagIds:     []string{criteria.Name},
		}
	}
	return DagRunRequest{
		OrderBy:          "execution_date",
		PageOffset:       0,
		PageLimit:        pageLimit,
		DagIds:           []string{criteria.Name},
		ExecutionDateGte: criteria.StartDate.Format(airflowDateFormat),
		ExecutionDateLte: criteria.EndDate.Format(airflowDateFormat),
	}
}

func (s *Airflow) Clear(ctx context.Context, jobName string, executionTime time.Time) error {
	return s.ClearBatch(ctx, jobName, executionTime, executionTime)
}

func (s *Airflow) ClearBatch(ctx context.Context, jobName string, startExecutionTime, endExecutionTime time.Time) error {
	data := []byte(fmt.Sprintf(`{"start_date": %q, "end_date": %q, "dry_run": false, "reset_dag_runs": true, "only_failed": false}`,
		startExecutionTime.UTC().Format(airflowDateFormat),
		endExecutionTime.UTC().Format(airflowDateFormat)))
	req := Request{
		Path:   fmt.Sprintf(dagRunClearURL, jobName),
		Method: http.MethodPost,
		Body:   data,
	}

	_, err := s.client.Invoke(ctx, req, s.auth)
	if err != nil {
		return fmt.Errorf("failure while clearing airflow dag runs: %w", err)
	}
	return nil
}

func (s *Airflow) CancelRun(ctx context.Context, jobName string, dagRunID string) error {
	data := []byte(`{"state": "failed"}`)
	req := Request{
		Path:   fmt.Sprintf(dagRunModifyURL, jobName, dagRunID),
		Method: http.MethodPatch,
		Body:   data,
	}

	_, err := s.client.Invoke(ctx, req, s.auth)
	if err != nil {
		return fmt.Errorf("failure while canceling airflow dag run: %w", err)
	}
	return nil
}

func (s *Airflow) CreateRun(ctx context.Context, jobName string, executionTime time.Time, dagRunIDPrefix string) error {
	data := []byte(fmt.Sprintf(`{"dag_run_id": %q, "execution_date": %q}`,
		fmt.Sprintf("%s__%s", dagRunIDPrefix,
			executionTime.UTC().Format(airflowDateFormat)),
		executionTime.UTC().Format(airflowDateFormat)),
	)

	req := Request{
		Path:   fmt.Sprintf(dagRunCreateURL, jobName),
		Method: http.MethodPost,
		Body:   data,
	}

	_, err := s.client.Invoke(ctx, req, s.auth)
	if err != nil {
		return fmt.Errorf("failure while creating airflow dag run: %w", err)
	}
	return nil
}

func (s *Airflow) TaskInstances(ctx context.Context, dagID string, dagRunID string) (*TaskInstances, error) {
	req := Request{
		Path:   fmt.Sprintf(taskInstances, dagID, dagRunID),
		Method: http.MethodGet,
	}

	resp, err := s.client.Invoke(ctx, req, s.auth)
	if err != nil {
		return nil, fmt.Errorf("unable to invoke task instances: %w", err)
	}

	var tasks TaskInstances
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("json error on parsing airflow task runs: %s, %w", string(resp), err)
	}

	return &tasks, nil
}

func NewAirflow(auth Auth) *Airflow {
	client := NewAirflowClient()
	return &Airflow{
		client: client,
		auth:   auth,
	}
}
