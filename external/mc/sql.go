package mc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aliyun/aliyun-odps-go-sdk/sqldriver"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

func NewSQLClientFromConfig(cfg *config.Config) (*sql.DB, error) {
	mcCreds := os.Getenv(MaxcomputeAccount)
	if mcCreds != "" {
		return NewSQLClient(mcCreds)
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

	return NewSQLClient(acc)
}

func NewSQLClient(creds string) (*sql.DB, error) {
	var c1 maxComputeCredentials
	if err := json.Unmarshal([]byte(creds), &c1); err != nil {
		return nil, err
	}

	conf := sqldriver.Config{
		AccessId:             c1.AccessID,
		AccessKey:            c1.AccessKey,
		StsToken:             "",
		Endpoint:             c1.Endpoint,
		ProjectName:          c1.ProjectName,
		HttpTimeout:          0,
		TcpConnectionTimeout: 30 * time.Second,
	}

	dsn := conf.FormatDsn()
	db, err := sql.Open("odps", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
