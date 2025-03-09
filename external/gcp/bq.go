package gcp

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

const (
	BigqueryAccount = "BIGQUERY_ACCOUNT"

	connTimeout = 5 * time.Second
)

//scopes := []string{
//	"https://www.googleapis.com/auth/drive",
//	"https://www.googleapis.com/auth/drive.file",
//	"https://www.googleapis.com/auth/drive.readonly",
//	"https://www.googleapis.com/auth/spreadsheets",
//	"https://www.googleapis.com/auth/spreadsheets.readonly",
//	bigquery.Scope,
//}

type ClientProvider struct {
	bq         *bigquery.Client
	drive      *drive.Service
	sheets     *sheets.Service
	static     bool
	staticCred string

	profile   *config.Profile
	clientMap map[string]*bigquery.Client
	driveMap  map[string]*drive.Service
	sheetMap  map[string]*sheets.Service
	mu        sync.Mutex
}

func NewClientProvider(cfg *config.Config) (*ClientProvider, error) {
	profile := cfg.GetCurrentProfile()
	acc := os.Getenv(BigqueryAccount)
	if acc != "" {
		return &ClientProvider{
			staticCred: acc,
			static:     true,
			profile:    profile,
		}, nil
	}

	if !profile.Dynamic {
		key := profile.GCPCred
		if key == "" {
			return nil, errors.New("key not found for Bigquery account")
		}

		acc, err := keyring.Get(key)
		if err != nil {
			return nil, err
		}

		if acc != "" {
			return &ClientProvider{
				staticCred: acc,
				static:     true,
				profile:    profile,
			}, nil
		} else {
			return nil, errors.New("empty value for account")
		}
	}

	return &ClientProvider{
		static:    false,
		profile:   profile,
		clientMap: make(map[string]*bigquery.Client),
		driveMap:  make(map[string]*drive.Service),
		sheetMap:  make(map[string]*sheets.Service),
	}, nil
}

func (p *ClientProvider) GetClient(proj string, scopes ...string) (*bigquery.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.static {
		if p.bq != nil {
			return p.bq, nil
		}

		client, err := NewClient(p.staticCred, proj, scopes)
		if err != nil {
			return nil, err
		}
		p.bq = client
		return client, err
	}

	key := proj + "_gcp"
	client, ok := p.clientMap[key]
	if ok {
		return client, nil
	}

	credKey, err := p.profile.GetCred(key)
	if err != nil {
		return nil, err
	}

	client, ok = p.clientMap[credKey]
	if ok {
		p.clientMap[key] = client
		return client, nil
	}

	cred, err := keyring.Get(credKey)
	if err != nil {
		return nil, err
	}

	clnt, err := NewClient(cred, proj, scopes)
	if err != nil {
		return nil, err
	}

	p.clientMap[key] = clnt
	p.clientMap[credKey] = clnt
	return clnt, nil
}

func NewClient(svcAccount string, proj string, scopes []string) (*bigquery.Client, error) {
	scopes = append(scopes, bigquery.Scope)

	ctx, cancelfunc := context.WithTimeout(context.Background(), connTimeout)
	defer cancelfunc()

	cred, err := google.CredentialsFromJSON(ctx, []byte(svcAccount), scopes...)
	if err != nil {
		return nil, errors.New("failed to read account")
	}

	if proj == "" {
		proj = cred.ProjectID
	}

	c, err := bigquery.NewClient(ctx, proj, option.WithCredentials(cred))
	if err != nil {
		return nil, errors.New("failed to create BQ client")
	}

	return c, nil
}
