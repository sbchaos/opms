package gcp

import (
	"context"
	"errors"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

const (
	BigqueryAccount = "BIGQUERY_ACCOUNT"

	connTimeout = 5 * time.Second
)

type Client struct {
	*bigquery.Client
}

func NewClientFromConfig(cfg *config.Config) (*Client, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), connTimeout)
	defer cancelfunc()

	bqAccount := os.Getenv(BigqueryAccount)
	if a := os.Getenv(bqAccount); a != "" {
		return NewClient(ctx, bqAccount)
	}

	profile := cfg.GetCurrentProfile()
	key := profile.GCPCred
	if key == "" {
		return nil, errors.New("key not found for Bigquery account")
	}

	acc, err := keyring.Get(key)
	if err != nil {
		return nil, err
	}

	return NewClient(ctx, acc)
}

func NewClient(ctx context.Context, svcAccount string) (*Client, error) {
	scopes := []string{
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.file",
		"https://www.googleapis.com/auth/drive.readonly",
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/spreadsheets.readonly",
		bigquery.Scope,
	}
	cred, err := google.CredentialsFromJSON(ctx, []byte(svcAccount), scopes...)
	if err != nil {
		return nil, errors.New("failed to read account")
	}

	c, err := bigquery.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, errors.New("failed to create BQ client")
	}

	return &Client{c}, nil
}
