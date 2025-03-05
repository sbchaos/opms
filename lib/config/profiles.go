package config

import (
	"errors"

	"github.com/sbchaos/opms/lib/util"
)

type Profile struct {
	Name string `json:"name"`
	// The actual credentials are stored in the keyring
	MCCred  string `json:"mc,omitempty"`
	GCPCred string `json:"gcp,omitempty"`
	Airflow string `json:"airflow,omitempty"`

	Dynamic   bool              `json:"dynamic"`
	Creds     map[string]string `json:"creds,omitempty"`
	Variables map[string]any    `json:"variables,omitempty"`
}

func (p *Profile) Merge(other *Profile) {
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

	newVariables := util.MergeAnyMaps(other.Variables, p.Variables)
	p.Variables = newVariables
}

func (p *Profile) SetVariable(key string, value any) {
	p.Variables[key] = value
}

func (p *Profile) GetVariable(key string) (any, error) {
	val, ok := p.Variables[key]
	if !ok {
		return nil, errors.New("value not found for key: " + key)
	}
	return val, nil
}

func (p *Profile) SetCred(key string, val string) {
	p.Creds[key] = val
}

func (p *Profile) GetCred(key string) (string, error) {
	val, ok := p.Creds[key]
	if !ok {
		return "", errors.New("cred not found for key: " + key)
	}
	return val, nil
}
