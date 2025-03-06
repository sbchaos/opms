package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/sbchaos/opms/lib/keyring"
)

func (p *ClientProvider) GetDriveClient(proj string) (*drive.Service, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.static {
		if p.drive != nil {
			return p.drive, nil
		}

		client, err := NewDriveClient(p.staticCred)
		if err != nil {
			return nil, err
		}
		p.drive = client
		return client, err
	}

	key := proj + "_gcp"
	client, ok := p.driveMap[key]
	if ok {
		return client, nil
	}

	credKey, err := p.profile.GetCred(key)
	if err != nil {
		return nil, err
	}

	client, ok = p.driveMap[credKey]
	if ok {
		p.driveMap[key] = client
		return client, nil
	}

	cred, err := keyring.Get(credKey)
	if err != nil {
		return nil, err
	}

	clnt, err := NewDriveClient(cred)
	if err != nil {
		return nil, err
	}

	p.driveMap[key] = clnt
	p.driveMap[credKey] = clnt
	return clnt, nil
}

func NewDriveClient(svcAccount string) (*drive.Service, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), connTimeout)
	defer cancelfunc()

	srv, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(svcAccount)))
	if err != nil {
		return nil, fmt.Errorf("not able to create drive service err: %w", err)
	}

	return srv, nil
}
