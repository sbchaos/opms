// nolint
package macro

import (
	"time"
)

const (
	// taskConfigPrefix will be used to prefix all the config variables of
	// transformation instance, i.e. task
	taskConfigPrefix = "TASK__"

	// projectConfigPrefix will be used to prefix all the config variables of
	// a project, i.e. registered entities
	projectConfigPrefix = "GLOBAL__"

	contextProject       = "proj"
	contextSecret        = "secret"
	contextSystemDefined = "inst"
	contextTask          = "task"

	TimeISOFormat = time.RFC3339
	TimeSQLFormat = time.DateTime

	// Configuration for system defined variables
	configDstart        = "DSTART"
	configStartDate     = "START_DATE"
	configDend          = "DEND"
	configEndDate       = "END_DATE"
	configExecutionTime = "EXECUTION_TIME"
	configScheduleTime  = "SCHEDULE_TIME"
	configScheduleDate  = "SCHEDULE_DATE"
	configDestination   = "JOB_DESTINATION"

	JobAttributionLabelsKey = "JOB_LABELS"
)

func getTimeConfigs(start, end, executedAt, scheduledAt time.Time) map[string]any {
	vars := map[string]any{
		configDstart:        start.Format(TimeISOFormat),
		configDend:          end.Format(TimeISOFormat),
		configExecutionTime: executedAt.Format(TimeSQLFormat),

		configStartDate:    start.Format(time.DateOnly),
		configEndDate:      end.Format(time.DateOnly),
		configScheduleTime: scheduledAt.Format(TimeISOFormat),
		configScheduleDate: scheduledAt.Format(time.DateOnly),
	}

	return vars
}
