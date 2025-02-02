package config

import "github.com/sbchaos/opms/lib/util"

type Profiles struct {
	Name string `json:"name"`
	// The actual credentials are stored in the keyring
	MCCred string `json:"mc,omitempty"`

	GCPCred string `json:"gcp,omitempty"`

	Airflow string `json:"airflow,omitempty"`

	// Creds can be used to store, alternate credentials
	Creds     map[string]string `json:"creds,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

func (p *Profiles) Merge(other *Profiles) {
	if p.MCCred == "" {
		p.MCCred = other.MCCred
	}
	if p.GCPCred == "" {
		p.GCPCred = other.GCPCred
	}
	if p.Airflow == "" {
		p.Airflow = other.Airflow
	}

	newCreds := util.MergeMaps(other.Creds, p.Creds)
	p.Creds = newCreds

	newVariables := util.MergeMaps(other.Variables, p.Variables)
	p.Variables = newVariables
}
