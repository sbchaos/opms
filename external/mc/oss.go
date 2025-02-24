package mc

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

func NewOSSClientFromConfig(cfg *config.Config) (*oss.Client, error) {
	mcCreds := os.Getenv(MaxcomputeAccount)
	if mcCreds != "" {
		return NewOssClient(mcCreds)
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

	return NewOssClient(acc)
}

func NewOssClient(creds string) (*oss.Client, error) {
	var c1 maxComputeCredentials
	if err := json.Unmarshal([]byte(creds), &c1); err != nil {
		return nil, err
	}

	credProvider := credentials.NewStaticCredentialsProvider(c1.AccessID, c1.AccessKey, c1.SecurityToken)
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithEndpoint(c1.OSSEndpoint).
		WithRegion(c1.Region)

	if cfg.CredentialsProvider == nil {
		return nil, errors.New("OSS: credentials provider is required")
	}

	return oss.NewClient(cfg), nil
}
