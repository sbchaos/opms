package mc

import (
	"encoding/json"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
)

type maxComputeCredentials struct {
	AccessID    string `json:"access_id"`
	AccessKey   string `json:"access_key"`
	Endpoint    string `json:"mc_endpoint"`
	ProjectName string `json:"project_name"`
}

func NewClient(creds string) (*odps.Odps, error) {
	var c1 maxComputeCredentials
	if err := json.Unmarshal([]byte(creds), &c1); err != nil {
		return nil, err
	}

	aliAccount := account.NewAliyunAccount(c1.AccessID, c1.AccessKey)
	odpsIns := odps.NewOdps(aliAccount, c1.Endpoint)

	return odpsIns, nil
}
