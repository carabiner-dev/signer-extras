// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spiffe

import (
	"context"
	"fmt"
	"os"

	"github.com/sigstore/sigstore/pkg/oauthflow"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"

	"github.com/carabiner-dev/signer/sts"
)

var _ sts.Provider = &Spire{}

type Spire struct{}

const (
	defaultSocketPath = "/var/run/spire/agent.sock"
)

// Register registers the spire provider in the sts DefaultProviders collection
func Register() {
	sts.RegisterProvider("spire", &Spire{})
}

// Provide returns a token for the spire agent by fetching a SVID from the
// socket. The socket location is read from SPIFFE_ENDPOINT_SOCKET or from
// /var/run/spire/agent.sock if unset.
func (spire *Spire) Provide(ctx context.Context, audience string) (*oauthflow.OIDCIDToken, error) {
	path := defaultSocketPath
	if p := os.Getenv("SPIFFE_ENDPOINT_SOCKET"); p != "" {
		path = p
	}

	// If the path is not found, then we asume we're not in a spiffe env
	if _, err := os.Stat(path); err != nil {
		return nil, nil
	}

	client, err := workloadapi.New(ctx, workloadapi.WithAddr(fmt.Sprintf("unix://%s", path)))
	if err != nil {
		return nil, fmt.Errorf("building spire client: %w", err)
	}
	defer client.Close()

	// Fetch the SVID from the spire server
	svid, err := client.FetchJWTSVID(ctx, jwtsvid.Params{
		Audience: audience,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching SVID: %w", err)
	}

	// Extract the subject claim
	var subject string
	if sub, ok := svid.Claims["sub"]; ok {
		if s, ok := sub.(string); ok {
			subject = s
		}
	}

	// Return the token in the sigstore wrapper
	return &oauthflow.OIDCIDToken{
		RawString: svid.Marshal(),
		Subject:   subject,
	}, nil
}
