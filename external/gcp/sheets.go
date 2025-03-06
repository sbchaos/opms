package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/sbchaos/opms/lib/keyring"
)

func (p *ClientProvider) GetSheetsClient(proj string) (*sheets.Service, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.static {
		if p.sheets != nil {
			return p.sheets, nil
		}

		client, err := NewSheetsClient(p.staticCred)
		if err != nil {
			return nil, err
		}
		p.sheets = client
		return client, err
	}

	key := proj + "_gcp"
	client, ok := p.sheetMap[key]
	if ok {
		return client, nil
	}

	credKey, err := p.profile.GetCred(key)
	if err != nil {
		return nil, err
	}

	client, ok = p.sheetMap[credKey]
	if ok {
		p.sheetMap[key] = client
		return client, nil
	}

	cred, err := keyring.Get(credKey)
	if err != nil {
		return nil, err
	}

	clnt, err := NewSheetsClient(cred)
	if err != nil {
		return nil, err
	}

	p.sheetMap[key] = clnt
	p.sheetMap[credKey] = clnt
	return clnt, nil
}

func NewSheetsClient(svcAccount string) (*sheets.Service, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), connTimeout)
	defer cancelfunc()

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(svcAccount)))
	if err != nil {
		return nil, fmt.Errorf("not able to create sheets service err: %w", err)
	}

	return srv, nil
}
