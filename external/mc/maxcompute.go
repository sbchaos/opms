package mc

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

const (
	MaxcomputeAccount = "MAXCOMPUTE_ACCOUNT"
)

type maxComputeCredentials struct {
	AccessID      string `json:"access_id"`
	AccessKey     string `json:"access_key"`
	McEndpoint    string `json:"mc_endpoint"`
	ProjectName   string `json:"project_name"`
	OSSEndpoint   string `json:"oss_endpoint"`
	Region        string `json:"region"`
	SecurityToken string `json:"security_token"`
}

func NewClientFromConfig(cfg *config.Config) (*odps.Odps, error) {
	mcCreds := os.Getenv(MaxcomputeAccount)
	if mcCreds != "" {
		return NewClient(mcCreds)
	}

	profile := cfg.GetCurrentProfile()
	key := profile.MCCred
	if key == "" {
		return nil, errors.New("key not found for Maxcompute account")
	}

	acc, err := keyring.Get(key)
	if err != nil {
		return nil, err
	}

	return NewClient(acc)
}

func NewClient(creds string) (*odps.Odps, error) {
	var c1 maxComputeCredentials
	if err := json.Unmarshal([]byte(creds), &c1); err != nil {
		return nil, err
	}

	aliAccount := account.NewAliyunAccount(c1.AccessID, c1.AccessKey)
	odpsIns := odps.NewOdps(aliAccount, c1.McEndpoint)
	odpsIns.SetDefaultProjectName(c1.ProjectName)

	return odpsIns, nil
}
