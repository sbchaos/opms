package plan

type Plan struct {
	ProjectName string                 `json:"project_name"`
	Job         OpsPerNS[JobPlan]      `json:"job"`
	Resource    OpsPerNS[ResourcePlan] `json:"resource"`
}

type OpsPerNS[T JobPlan | ResourcePlan] struct {
	Create  map[string][]T `json:"create"`
	Delete  map[string][]T `json:"delete"`
	Update  map[string][]T `json:"update"`
	Migrate map[string][]T `json:"migrate"`
}

type JobPlan struct {
	Name         string  `json:"name"`
	OldNamespace *string `json:"old_namespace"`
	Path         string  `json:"path"`
}

type ResourcePlan struct {
	Name         string  `json:"name"`
	Datastore    string  `json:"datastore"`
	OldNamespace *string `json:"old_namespace"`
	Path         string  `json:"path"`
}
